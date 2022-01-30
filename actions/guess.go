package actions

import (
	"guess_my_word/internal/model"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type guess struct {
	Guesses int       `form:"guesses"`
	Word    string    `form:"word"`
	Mode    string    `form:"mode"`
	Start   time.Time `form:"start" time_format:"unix"`
	TZ      int       `form:"tz"`
}

type guessReply struct {
	Guess   string `json:"guess"`
	Correct bool   `json:"correct"`
	After   bool   `json:"after"`
	Before  bool   `json:"before"`
	Error   string `json:"error,omitempty"`
}

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

// GuessHandler is an API handler to process a user's guess.
func GuessHandler(c *gin.Context) {
	guess := guess{}
	reply := guessReply{}

	// Validate the guess
	if err := c.ShouldBind(&guess); err != nil {
		log.Println("Invalid request received: ", err)
		reply.Error = ErrInvalidRequest
	} else if len(strings.TrimSpace(guess.Word)) == 0 {
		reply.Error = ErrEmptyGuess
	} else if !wordStore.Validate(c, guess.Word) {
		reply.Error = ErrInvalidWord
	} else if guess.Start.Unix() == 0 {
		reply.Error = ErrInvalidStartTime
	}
	reply.Guess = strings.TrimSpace(guess.Word)
	if reply.Error != "" {
		c.JSON(200, reply)
		return
	}

	// Generate the word for the day
	tm := convertUTCToUser(guess.Start, guess.TZ)
	word, err := wordStore.GetForDay(c, tm, guess.Mode)
	if err != nil {
		reply.Error = err.Error()
		c.JSON(500, reply)
		return
	}

	if reply.Error == "" {
		switch strings.Compare(reply.Guess, word.Value) {
		case -1:
			reply.After = true
		case 1:
			reply.Before = true
		case 0:
			reply.Correct = true
		}
	}

	if reply.Correct {
		word.Guesses = append(word.Guesses, model.Guess{
			Count: guess.Guesses,
		})

		// TODO: Centralize this key logic
		key := guess.Mode + "/day/" + tm.Format("2006-01-02")
		dataStore.SetWord(c, key, word)
	}

	c.JSON(200, reply)
}
