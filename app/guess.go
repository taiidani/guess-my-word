package app

import (
	"context"
	"fmt"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"guess_my_word/internal/sessions"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// ErrInvalidWord is emitted when the guess is not a legitimate word
	ErrInvalidWord = "Guess must be a valid Scrabble word"

	// ErrInvalidStartTime is emitted when the start time is malformed or invalid
	ErrInvalidStartTime = "Invalid start time provided with request"

	// ErrEmptyGuess is emitted when the guess provided was empty
	ErrEmptyGuess = "Guess must not be empty"
)

var guessMutex = sync.Mutex{}

// GuessHandler is an API handler to process a user's guess.
func GuessHandler(c *gin.Context) {
	session, err := startSession(c)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	word := strings.TrimSpace(c.Request.PostFormValue("word"))

	// Validate the guess
	if len(word) == 0 {
		c.HTML(http.StatusBadRequest, "error.gohtml", ErrEmptyGuess)
		return
	} else if !wordStore.Validate(c, word) {
		c.HTML(http.StatusBadRequest, "error.gohtml", ErrInvalidWord)
		return
	}

	err = guessHandlerReply(c, session, word)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	if err := session.Save(); err != nil {
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	c.HTML(http.StatusOK, "guesser.gohtml", session.Current())
}

var fnSetEndTime func() *time.Time = func() *time.Time {
	now := time.Now()
	return &now
}

func guessHandlerReply(ctx context.Context, session *sessions.Session, guess string) error {
	// Only one guess operation may happen simultaneously
	// This allows us to get the Word then modify it with new data without overriding anyone
	// else's contributions
	guessMutex.Lock()
	defer guessMutex.Unlock()

	// Generate the word for the day
	tm := session.DateUser()
	word, err := wordStore.GetForDay(ctx, tm, session.Mode)
	if err != nil {
		return err
	}

	current := session.Current()
	switch strings.Compare(guess, word.Value) {
	case -1:
		for _, w := range current.Before {
			if w == guess {
				return fmt.Errorf("you have already guessed this word")
			}
		}
		current.Before = append(current.Before, guess)
		sort.Strings(current.Before)
	case 1:
		for _, w := range current.After {
			if w == guess {
				return fmt.Errorf("you have already guessed this word")
			}
		}
		current.After = append(current.After, guess)
		sort.Strings(current.After)
	case 0:
		current.Answer = guess
		current.End = fnSetEndTime()

		// Record the successful guess
		word.Guesses = append(word.Guesses, model.Guess{
			Count: len(current.Before) + len(current.After) + 1,
		})

		return wordStore.SetWord(ctx, datastore.WordKey(session.Mode, tm), word)
	}

	return nil
}
