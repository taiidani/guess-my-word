package app

import (
	"fmt"
	"guess_my_word/internal/sessions"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type indexBag struct {
	baseBag
	Mode    string
	History *sessions.SessionMode
}

func IndexHandler(c *gin.Context) {
	data := indexBag{}
	data.Page = "home"

	s := sessions.New(c)
	defer func() {
		if err := s.Save(); err != nil {
			slog.Warn("Unable to save session", "error", err)
		}
	}()

	// Assign the mode
	switch {
	// Allow mode selection through the path
	case c.Param("mode") != "":
		data.Mode = c.Param("mode")
	default:
		data.Mode = "default"
	}
	s.Mode = data.Mode

	// Load list data for the current mode
	// This also validates that it is an existing mode
	var err error
	data.List, err = listStore.GetList(c, data.Mode)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.gohtml", fmt.Sprintf("Could not load list %q: %s", data.Mode, err))
		return
	}

	data.History = s.Current()

	c.HTML(http.StatusOK, "index.gohtml", data)
}
