package actions

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func middlewareStandardHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		csp := []string{
			`default-src 'self' cdn.jsdelivr.net`,
			`script-src 'self' *.jquery.com cdn.jsdelivr.net 'unsafe-eval'`,
			`style-src 'self' cdn.jsdelivr.net`,
			`frame-ancestors 'none'`,
			`form-action 'self'`,
			`base-uri 'self'`,
			`img-src 'self' data:`,
		}
		c.Header("Content-Security-Policy", strings.Join(csp, "; "))

		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "no-referrer")

		if gin.Mode() == gin.ReleaseMode {
			c.Header("Strict-Transport-Security", "max-age=63072000") // 2 years
		}
	}
}
