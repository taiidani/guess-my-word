package app

import (
	"fmt"
	"guess_my_word/internal/model"
	"guess_my_word/internal/sessions"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	data := struct {
		Mode    string
		List    model.List
		History *sessions.SessionMode
		Lists   []string
	}{}

	s := sessions.New(c)
	defer func() {
		if err := s.Save(); err != nil {
			log.Printf("WARN: Unable to save session: %s", err)
		}
	}()

	// Load the lists
	var err error
	data.Lists, err = listStore.GetLists(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", "Unable to load word lists")
		return
	}

	// Allow mode selection through the select dropdown
	data.Mode = strings.ToLower(c.Request.URL.Query().Get("mode"))
	if data.Mode == "" {
		data.Mode = "default"
	}

	// Load list data for the current mode
	// This also validates that it is an existing mode
	data.List, err = listStore.GetList(c, data.Mode)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.gohtml", fmt.Sprintf("Could not load list %q: %s", data.Mode, err))
		return
	}

	s.Mode = data.Mode
	data.History = s.Current()

	c.HTML(http.StatusOK, "index.gohtml", data)
}
