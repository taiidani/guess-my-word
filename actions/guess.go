package actions

import (
	"context"
	"fmt"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"log"
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

	// ErrInvalidTimezone is emitted when the timezone is malformed or invalid
	ErrInvalidTimezone = "Invalid timezone provided with request"

	// ErrEmptyGuess is emitted when the guess provided was empty
	ErrEmptyGuess = "Guess must not be empty"
)

var guessMutex = sync.Mutex{}

// GuessHandler is an API handler to process a user's guess.
func GuessHandler(c *gin.Context) {
	request, err := parseBodyData(c)
	if err != nil {
		log.Println("Unable to parse body data: ", err)
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

	err = guessHandlerReply(c, request, word)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	if err := request.Session.Save(); err != nil {
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	c.HTML(http.StatusOK, "guesser.gohtml", request.Session.Current())
}

func guessHandlerReply(ctx context.Context, data bodyData, guess string) error {
	// Only one guess operation may happen simultaneously
	// This allows us to get the Word then modify it with new data without overriding anyone
	// else's contributions
	guessMutex.Lock()
	defer guessMutex.Unlock()

	// Generate the word for the day
	tm := data.Session.DateUser(data.TZ)
	word, err := wordStore.GetForDay(ctx, tm, data.Session.Mode)
	if err != nil {
		return err
	}

	current := data.Session.Current()
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
		now := time.Now()
		current.End = &now

		// Record the successful guess
		word.Guesses = append(word.Guesses, model.Guess{
			Count: len(current.Before) + len(current.After) + 1,
		})

		wordStore.SetWord(ctx, datastore.WordKey(data.Session.Mode, tm), word)
	}

	return nil
}
