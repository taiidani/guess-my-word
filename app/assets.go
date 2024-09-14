package app

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed assets
var assets embed.FS

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	if dev {
		slog.Debug("Serving file from disk", "path", r.URL.Path)

		// Sanitize the path
		cwd, _ := os.Getwd()
		safePath := filepath.Join(cwd, "app")
		absPath, err := filepath.Abs(filepath.Join(safePath, r.URL.Path))
		if err != nil || !strings.HasPrefix(absPath, safePath) {
			errorResponse(w, http.StatusBadRequest, fmt.Errorf("Invalid file name %q", absPath))
			return
		}

		http.ServeFile(w, r, filepath.Join("app", r.URL.Path))
	} else {
		slog.Debug("Serving file from FS", "path", r.URL.Path)
		http.ServeFileFS(w, r, assets, r.URL.Path)
	}
}
