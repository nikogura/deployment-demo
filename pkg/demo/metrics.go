package demo

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// RED metrics — Rate, Errors, Duration.

// RequestsTotal counts total HTTP requests by method, path, and status code.
//
//nolint:gochecknoglobals // Prometheus metrics are registered at package level.
var RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "demo_http_requests_total",
	Help: "Total HTTP requests by method, path, and status code.",
}, []string{"method", "path", "code"})

// RequestDuration observes HTTP request durations in seconds.
//
//nolint:gochecknoglobals // Prometheus metrics are registered at package level.
var RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "demo_http_request_duration_seconds",
	Help:    "HTTP request duration in seconds.",
	Buckets: prometheus.DefBuckets,
}, []string{"method", "path"})

// RequestErrors counts HTTP requests that resulted in 5xx status codes.
//
//nolint:gochecknoglobals // Prometheus metrics are registered at package level.
var RequestErrors = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "demo_http_request_errors_total",
	Help: "Total HTTP requests that resulted in 5xx errors.",
}, []string{"method", "path"})

// HealthCheckTotal counts health check results.
//
//nolint:gochecknoglobals // Prometheus metrics are registered at package level.
var HealthCheckTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "demo_health_checks_total",
	Help: "Total health check results by status.",
}, []string{"status"})
