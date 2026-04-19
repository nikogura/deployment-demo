package demo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Config holds runtime configuration from environment variables.
type Config struct {
	Port        string // DEMO_PORT, default ":8080"
	TempoURL    string // DEMO_TEMPO_URL, empty = tracing disabled
	ServiceName string // DEMO_SERVICE_NAME, default "deployment-demo"
}

// LoadConfig reads configuration from environment variables.
func LoadConfig() (cfg Config) {
	cfg.Port = envOrDefault("DEMO_PORT", ":8080")
	cfg.TempoURL = os.Getenv("DEMO_TEMPO_URL")
	cfg.ServiceName = envOrDefault("DEMO_SERVICE_NAME", "deployment-demo")
	return cfg
}

// Run starts the HTTP server. If Health=="broken", the server starts
// normally (so the UI, metrics, and health endpoint are all reachable)
// but /healthz returns 503. K8s liveness probes will eventually kill
// the pod, causing a crashloop that's visible in metrics and logs.
func Run(ctx context.Context, cfg Config, logger *slog.Logger) (err error) {
	// Initialize tracing.
	var tracerShutdown func(context.Context) error
	tracerShutdown, err = InitTracer(ctx, cfg.ServiceName, cfg.TempoURL, logger)
	if err != nil {
		err = fmt.Errorf("init tracer: %w", err)
		return err
	}
	defer func() { _ = tracerShutdown(ctx) }()

	// Set up routes.
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	SetupUIRoutes(mux)

	// Wrap with instrumentation middleware (metrics + tracing + logging).
	handler := Instrument(mux, logger)

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	logger.InfoContext(ctx, "server starting",
		"address", cfg.Port,
		"version", Version,
		"theme", Theme,
		"health", Health,
		"service", cfg.ServiceName,
	)

	if !IsHealthy() {
		logger.WarnContext(ctx, "THIS BUILD IS INTENTIONALLY BROKEN — /healthz will return 503",
			"version", Version,
			"health", Health,
		)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.InfoContext(ctx, "shutting down", "version", Version)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = srv.Shutdown(shutdownCtx)
		return err
	case listenErr := <-errCh:
		err = listenErr
		return err
	}
}

func envOrDefault(key, fallback string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}
