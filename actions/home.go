package actions

import (
	"guess_my_word/internal/words"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// HomeHandler displays the home page HTML
func HomeHandler(c *gin.Context) {
	mode := strings.ToLower(strings.TrimSpace(c.Query("mode")))
	if mode == "" {
		mode = "default"
	}

	tm := time.Now().UTC().AddDate(0, 0, -1)
	yesterday, _ := words.GetForDay(c, tm, mode)

	c.HTML(200, "index.html", gin.H{
		"debug":     gin.IsDebugging(),
		"mode":      mode,
		"yesterday": yesterday,
	})
}
