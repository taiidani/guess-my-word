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
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	if err := session.Clear(); err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", err.Error())
		return
	}

	c.Redirect(301, "/")
}
