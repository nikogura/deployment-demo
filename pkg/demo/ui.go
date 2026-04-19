package demo

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/nikogura/deployment-demo/pkg/ui"
)

// SetupUIRoutes serves the embedded SPA with fallback to index.html.
func SetupUIRoutes(mux *http.ServeMux) {
	uiFS, err := fs.Sub(ui.Files, "dist")
	if err != nil {
		mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{
				"message": "deployment-demo API running. UI not built — run 'make build-ui' first.",
			})
		})
		return
	}

	fileServer := http.FileServer(http.FS(uiFS))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/api/") || path == "/healthz" || path == "/readyz" || path == "/metrics" {
			http.NotFound(w, r)
			return
		}

		if path != "/" {
			cleanPath := strings.TrimPrefix(path, "/")
			_, openErr := fs.Stat(uiFS, cleanPath)
			if openErr == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		indexContent, readErr := fs.ReadFile(uiFS, "index.html")
		if readErr != nil {
			http.Error(w, "UI not available", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		_, _ = w.Write(indexContent)
	})
}
