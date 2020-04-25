package actions

import (
	"time"

	"github.com/gin-gonic/gin"
)

// HomeHandler displays the home page HTML
func HomeHandler(c *gin.Context) {
	yesterday, _ := generateWord(time.Now().UTC().AddDate(0, 0, -1))

	c.HTML(200, "index.html", gin.H{
		"yesterday": yesterday,
	})
}
