package actions

import (
	"io"
	"io/fs"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func newStaticHandler(assets fs.FS) func(c *gin.Context) {
	return func(c *gin.Context) {
		staticHandler(c, assets)
	}
}

// staticHandler handles Static assets
func staticHandler(c *gin.Context, assets fs.FS) {
	p := c.Param("filepath")
	log.Println("/assets" + p)
	f, err := assets.Open("assets" + p)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	contentType := "text/plain"
	switch {
	case strings.HasSuffix(p, ".svg"):
		contentType = "image/svg+xml"
	case strings.HasSuffix(p, ".css"):
		contentType = "text/css"
	case strings.HasSuffix(p, ".js"):
		contentType = "application/javascript"
	case strings.HasSuffix(p, ".png"):
		contentType = "image/png"
	}

	c.Header("Content-Type", contentType)
	c.Writer.Write(body)
}
