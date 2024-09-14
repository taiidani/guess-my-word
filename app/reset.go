package app

import (
	"log/slog"
	"net/http"
)

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	session, err := startSession(w, r)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := session.Clear(); err != nil {
		errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	renderRedirect(w, 301, "/")
}
