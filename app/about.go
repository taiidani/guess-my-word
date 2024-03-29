package app

import (
	"guess_my_word/internal/sessions"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type aboutBag struct {
	baseBag
}

func AboutHandler(c *gin.Context) {
	data := aboutBag{}
	data.Page = "about"

	data.Session = sessions.New(c)
	defer func() {
		if err := data.Session.Save(); err != nil {
			slog.Warn("Unable to save session", "error", err)
		}
	}()

	c.HTML(http.StatusOK, "about.gohtml", data)
}
