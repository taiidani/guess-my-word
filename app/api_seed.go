package app

import (
	"guess_my_word/internal/sessions"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type seedReply struct {
	Error string `json:"error,omitempty"`
}

// SeedHandler is an internal API handler for pre-populating data to test with.
func SeedHandler(c *gin.Context) {
	session, err := startSession(c)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		c.JSON(http.StatusBadRequest, seedReply{Error: err.Error()})
		return
	}

	// Seed the session with predictable data
	session.Mode = "default"
	session.History = map[string]*sessions.SessionMode{
		// The answer is "website"
		// Yesterday's answer is "worst"
		"default": {
			Start:  time.Date(2022, 11, 7, 0, 0, 0, 0, time.UTC),
			Before: []string{},
			After:  []string{},
		},

		// The answer is "gemshorn"
		// Yesterday's answer is "gabbroid"
		"hard": {
			Start:  time.Date(2022, 11, 7, 0, 0, 0, 0, time.UTC),
			Before: []string{},
			After:  []string{},
		},
	}

	if err := session.Save(); err != nil {
		c.JSON(http.StatusBadRequest, seedReply{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, seedReply{})
}
