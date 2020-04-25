package actions

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
)

// StaticHandler handles Static assets
func StaticHandler(c *gin.Context) {
	p := c.Param("filepath")
	log.Println("/assets" + p)
	f, err := pkger.Open("/assets" + p)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)
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
