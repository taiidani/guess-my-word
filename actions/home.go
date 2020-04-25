package actions

import (
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
	yesterday, _ := generateWord(tm, getWordList(mode))

	c.HTML(200, "index.html", gin.H{
		"mode":      mode,
		"yesterday": yesterday,
	})
}
