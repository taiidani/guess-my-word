package app

import (
	_ "embed"
	"net/http"

	"guess_my_word/internal/sessions"

	"github.com/gin-gonic/gin"
)

func ModeSetHandler(c *gin.Context) {
	s := sessions.New(c)

	newMode := c.Request.PostFormValue("mode")
	if newMode == "" {
		c.HTML(http.StatusBadRequest, "error.gohtml", "Invalid mode specified")
		return
	}

	s.Mode = newMode
	if err := s.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
