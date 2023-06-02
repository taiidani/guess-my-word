package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResetHandler(c *gin.Context) {
	request, err := parseBodyData(c)
	if err != nil {
		log.Println("Unable to parse body data: ", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	if err := request.Session.Clear(); err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", err.Error())
		return
	}

	c.Redirect(301, "/")
}
