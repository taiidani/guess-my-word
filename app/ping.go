package app

import (
	"github.com/gin-gonic/gin"
)

// PingHandler interacts with service healthchecks
func PingHandler(c *gin.Context) {
	c.String(200, "%s", "pong")
}
