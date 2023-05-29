package actions

import (
	"guess_my_word/internal/sessions"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	type indexData struct {
		Mode    string
		History *sessions.SessionMode
		List    []string
	}
	s := sessions.New(c)
	defer s.Save()

	// Load the lists
	lists, err := listStore.GetLists(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.gohtml", "Unable to load word lists")
		return
	}

	// Allow mode selection through the select dropdown
	mode := strings.ToLower(c.Request.URL.Query().Get("mode"))
	if mode == "" {
		mode = "default"
	}

	// Validate that the mode is real before setting it in the session
	var found bool
	for _, l := range lists {
		if mode == l {
			found = true
			break
		}
	}
	if !found {
		c.HTML(http.StatusBadRequest, "error.gohtml", "Invalid mode specified")
		return
	}
	s.Mode = mode

	data := indexData{
		Mode:    s.Mode,
		History: s.Current(),
		List:    lists,
	}

	c.HTML(http.StatusOK, "index.gohtml", data)
}
