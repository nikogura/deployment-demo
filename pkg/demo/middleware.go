package demo

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
)

// instrumentedWriter wraps http.ResponseWriter to capture the status code.
type instrumentedWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (w *instrumentedWriter) WriteHeader(code int) {
	if !w.written {
		w.statusCode = code
		w.written = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *instrumentedWriter) Write(b []byte) (n int, err error) {
	if !w.written {
		w.statusCode = http.StatusOK
		w.written = true
	}
	n, err = w.ResponseWriter.Write(b)
	return n, err
}

// Instrument wraps an http.Handler with RED metrics, tracing, and
// structured logging.
func Instrument(next http.Handler, logger *slog.Logger) (handler http.Handler) {
	tracer := otel.Tracer("deployment-demo")

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		method := r.Method

		ctx, span := tracer.Start(r.Context(), method+" "+path)
		defer span.End()

		iw := &instrumentedWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(iw, r.WithContext(ctx))

		duration := time.Since(start)
		code := strconv.Itoa(iw.statusCode)

		RequestsTotal.WithLabelValues(method, path, code).Inc()
		RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

		if iw.statusCode >= 500 {
			RequestErrors.WithLabelValues(method, path).Inc()
		}

		logger.InfoContext(ctx, "request",
			"method", method,
			"path", path,
			"status", iw.statusCode,
			"duration_ms", duration.Milliseconds(),
			"version", Version,
		)
	})
	return handler
}
