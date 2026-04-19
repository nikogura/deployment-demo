// Package demo implements the deployment-demo HTTP server.
package demo

// Build-time variables — set via ldflags to differentiate versions visually.
//
// Build with:
//
//	go build -ldflags "\
//	  -X github.com/nikogura/deployment-demo/pkg/demo.Version=1.0.0 \
//	  -X github.com/nikogura/deployment-demo/pkg/demo.Theme=green \
//	  -X github.com/nikogura/deployment-demo/pkg/demo.Health=ok"

// Version is the semver version string shown in the UI.
//
//nolint:gochecknoglobals // Must be var for ldflags injection.
var Version = "dev"

// Theme controls the visual color scheme of the SPA. Values: "green",
// "blue", "red". Each version gets a different theme so observers can
// instantly see which version is running.
//
//nolint:gochecknoglobals // Must be var for ldflags injection.
var Theme = "green"

// Health controls whether /healthz returns 200 or 503. Values: "ok"
// (healthy) or "broken" (returns 503, triggers K8s probe failure and
// eventual crashloop). The "broken" version still serves the UI and
// metrics so the failure is observable before the pod is killed.
//
//nolint:gochecknoglobals // Must be var for ldflags injection.
var Health = "ok"

// BuildTime is injected at build time.
//
//nolint:gochecknoglobals // Must be var for ldflags injection.
var BuildTime = "unknown"

// IsHealthy reports whether this build is configured as healthy.
func IsHealthy() (healthy bool) {
	healthy = Health == "ok"
	return healthy
}
