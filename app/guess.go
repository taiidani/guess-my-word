package app

import (
	"context"
	"errors"
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

type guessBag struct {
	Session *sessions.SessionMode

	Stats model.WordStats
}

var guessMutex = sync.Mutex{}

// GuessHandler is an API handler to process a user's guess.
func GuessHandler(c *gin.Context) {
	session, err := startSession(c)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	guess := strings.TrimSpace(c.Request.PostFormValue("word"))

	// Validate the guess
	if len(guess) == 0 {
		errorResponse(c, http.StatusBadRequest, errors.New(ErrEmptyGuess))
		return
	} else if !wordStore.Validate(c, guess) {
		errorResponse(c, http.StatusBadRequest, errors.New(ErrInvalidWord))
		return
	}

	word, err := guessHandlerReply(c, session, guess)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := session.Save(); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	data := fillGuessBag(session.Current(), wordStore, word)
	c.HTML(http.StatusOK, "guesser.gohtml", data)
}

func fillGuessBag(s *sessions.SessionMode, w wordClient, word model.Word) guessBag {
	bag := guessBag{}
	bag.Session = s
	bag.Stats = word.Stats()
	return bag
}

var fnSetEndTime func() *time.Time = func() *time.Time {
	now := time.Now()
	return &now
}

func guessHandlerReply(ctx context.Context, session *sessions.Session, guess string) (model.Word, error) {
	// Only one guess operation may happen simultaneously
	// This allows us to get the Word then modify it with new data without overriding anyone
	// else's contributions
	guessMutex.Lock()
	defer guessMutex.Unlock()

	// Generate the word for the day
	tm := session.DateUser()
	word, err := wordStore.GetForDay(ctx, tm, session.Mode)
	if err != nil {
		return word, err
	}

	current := session.Current()
	switch strings.Compare(guess, word.Value) {
	case -1:
		for _, w := range current.Before {
			if w == guess {
				return word, fmt.Errorf("you have already guessed this word")
			}
		}
		current.Before = append(current.Before, guess)
		sort.Strings(current.Before)
	case 1:
		for _, w := range current.After {
			if w == guess {
				return word, fmt.Errorf("you have already guessed this word")
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

		return word, wordStore.SetWord(ctx, datastore.WordKey(session.Mode, tm), word)
	}

	return word, nil
}
