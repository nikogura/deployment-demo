package demo

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitTracer sets up the OpenTelemetry trace provider. If tempoURL is
// empty, tracing is disabled (noop provider). Returns a shutdown func.
func InitTracer(ctx context.Context, serviceName, tempoURL string, logger *slog.Logger) (shutdown func(context.Context) error, err error) {
	if tempoURL == "" {
		logger.InfoContext(ctx, "tracing disabled (DEMO_TEMPO_URL not set)")
		shutdown = func(_ context.Context) (noopErr error) { return noopErr }
		return shutdown, err
	}

	var exporter sdktrace.SpanExporter
	exporter, err = otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(tempoURL),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		err = fmt.Errorf("create OTLP exporter: %w", err)
		return shutdown, err
	}

	var res *resource.Resource
	res, err = resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(Version),
		),
	)
	if err != nil {
		err = fmt.Errorf("create resource: %w", err)
		return shutdown, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	logger.InfoContext(ctx, "tracing enabled", "endpoint", tempoURL, "service", serviceName)
	shutdown = tp.Shutdown
	return shutdown, err
}
