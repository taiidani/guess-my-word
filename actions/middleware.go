package actions

import (
	"github.com/gin-gonic/gin"
)

func middlewareStandardHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		if gin.Mode() == gin.ReleaseMode {
			c.Header("Strict-Transport-Security", "max-age=63072000") // 2 years
		} else {
			c.Header("Access-Control-Allow-Origin", "*") // Let everyone in in dev mode only!
		}
	}
}
