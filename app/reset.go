package app

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResetHandler(c *gin.Context) {
	request, err := parseBodyData(c)
	if err != nil {
		slog.Warn("Unable to parse body data", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	if err := request.Session.Clear(); err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", err.Error())
		return
	}

	c.Redirect(301, "/")
}
