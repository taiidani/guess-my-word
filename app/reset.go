package app

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResetHandler(c *gin.Context) {
	session, err := startSession(c)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := session.Clear(); err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(301, "/")
}
