package actions

import (
	"strconv"
	"strings"
	"time"

	"guess_my_word/internal/data"

	"github.com/gobuffalo/buffalo"
)

type guess struct {
	word  string
	start time.Time
}

type guessReply struct {
	Guess   string `json:"guess"`
	Correct bool   `json:"correct"`
	After   bool   `json:"after"`
	Before  bool   `json:"before"`
	Error   string `json:"error"`
}

// GuessHandler is an API handler to process a user's guess.
func GuessHandler(c buffalo.Context) error {
	guess := extractGuess(c)
	reply := guessReply{}
	reply.Guess = guess.word

	// Validate the guess
	if len(reply.Guess) == 0 {
		reply.Error = "Guess must not be empty"
	} else if !data.ValidateWord(reply.Guess) {
		reply.Error = "Guess must be a valid Scrabble word"
	}

	// Load or generate the word for the day
	day, err := data.LoadDate(guess.start)
	if err != nil {
		// Couldn't load a day. Generate a new one
		c.Logger().Infof("Date object failed to load from datastore: %s", err)
		c.Logger().Info("Generating new date")
		day = data.NewDate(guess.start)
		if err := day.Save(); err != nil {
			return err
		}
	}

	if reply.Error == "" {
		switch strings.Compare(reply.Guess, day.Word) {
		case -1:
			reply.After = true
		case 1:
			reply.Before = true
		case 0:
			reply.Correct = true
		}
	}

	return c.Render(200, r.JSON(reply))
}

func extractGuess(c buffalo.Context) guess {
	ret := guess{}
	ret.word = strings.ToLower(strings.TrimSpace(c.Param("word")))

	startStr := strings.TrimSpace(c.Param("start"))
	if startUnix, err := strconv.ParseInt(startStr, 10, 64); err == nil {
		ret.start = time.Unix(startUnix/1000, 0)
	}

	return ret
}
