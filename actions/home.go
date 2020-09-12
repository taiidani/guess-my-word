package actions

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// HomeHandler displays the home page HTML
func HomeHandler(c *gin.Context) {
	mode := strings.ToLower(strings.TrimSpace(c.Query("mode")))
	if mode == "" {
		mode = "default"
	}

	c.HTML(200, "index.html", gin.H{
		"debug": gin.IsDebugging(),
		"mode":  mode,
	})
}
