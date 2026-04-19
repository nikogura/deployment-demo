package demo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//nolint:gochecknoglobals // Set once at startup, read-only thereafter.
var startTime = time.Now()

// RegisterRoutes registers all API routes on the mux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /readyz", handleHealth)
	mux.HandleFunc("GET /api/info", handleInfo)
	mux.Handle("GET /metrics", promhttp.Handler())
}

// handleHealth returns 200 if healthy, 503 if broken. The broken
// version still serves this endpoint (so you can see the 503 in the
// browser / metrics) — it just fails the K8s probe.
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	if IsHealthy() {
		HealthCheckTotal.WithLabelValues("healthy").Inc()
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "healthy",
			"version": Version,
		})
		return
	}

	HealthCheckTotal.WithLabelValues("unhealthy").Inc()
	writeJSON(w, http.StatusServiceUnavailable, map[string]string{
		"status":  "unhealthy",
		"version": Version,
		"reason":  "build configured as broken (Health=broken)",
	})
}

// InfoResponse is the payload for /api/info.
type InfoResponse struct {
	Version   string `json:"version"`
	Theme     string `json:"theme"`
	Health    string `json:"health"`
	Healthy   bool   `json:"healthy"`
	BuildTime string `json:"buildTime"`
	Uptime    string `json:"uptime"`
	UptimeSec int64  `json:"uptimeSec"`
}

// handleInfo returns build-time metadata for the SPA to render.
func handleInfo(w http.ResponseWriter, _ *http.Request) {
	uptime := time.Since(startTime)
	writeJSON(w, http.StatusOK, InfoResponse{
		Version:   Version,
		Theme:     Theme,
		Health:    Health,
		Healthy:   IsHealthy(),
		BuildTime: BuildTime,
		Uptime:    uptime.Round(time.Second).String(),
		UptimeSec: int64(uptime.Seconds()),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
