package actions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"guess_my_word/internal/data"

	"github.com/gobuffalo/buffalo"
)

type completeRequest struct {
	word       string        // The correct word, to validate that the user did guess correctly
	suggestion string        // The suggested word for tomorrow, if any
	username   string        // The user's username, if provided
	start      time.Time     // When the user started guessing. Determines which word they guessed
	duration   time.Duration // The amount of time that the user took
	guesses    uint64        // The number of guesses that the user took
}

type completeResponse struct {
	Error string `json:"error"`
}

// CompleteHandler is an API handler to process a user's leaderboard/suggestion entry.
func CompleteHandler(c buffalo.Context) error {
	request := extractCompleteRequest(c)
	response := completeResponse{}
	responseCode := 200

	// Validate the guess
	if len(request.word) == 0 {
		response.Error = "Answer must not be empty"
		responseCode = 400
	}

	// Load the date for this
	day, err := data.LoadDate(request.start)
	if err != nil {
		// Couldn't load a day. Generate a new one
		c.Logger().Infof("Date object failed to load from datastore: %s", err)
		response.Error = "Could not find entry for this date"
		responseCode = 404
	}

	// Only process if they have provided a username
	if response.Error == "" && len(request.username) > 0 {
		if request.word != day.Word {
			response.Error = "The given answer does not match the day provided"
			responseCode = 400
			return c.Render(responseCode, r.JSON(response))
		}

		c.Logger().Debugf("%#v", day)
		entry, err := day.AddLeaderboard(request.username, request.guesses, request.duration)
		if err != nil {
			response.Error = err.Error()
			responseCode = 400
			return c.Render(responseCode, r.JSON(response))
		}

		// Add a suggestion if the user has provided it
		if len(request.suggestion) > 0 {
			c.Logger().Debugf("User %s has suggested word %s", entry.User.Name, request.suggestion)
			day.AddSuggestion(request.suggestion, entry.User)
		}

		c.Logger().Debugf("%#v", day)
		if err := day.Save(); err != nil {
			response.Error = fmt.Sprintf("Unable to save day information: %s", err)
			responseCode = 400
		}
	}

	return c.Render(responseCode, r.JSON(response))
}

func extractCompleteRequest(c buffalo.Context) completeRequest {
	var err error
	ret := completeRequest{}
	ret.word = strings.ToLower(strings.TrimSpace(c.Param("word")))
	ret.suggestion = strings.ToLower(strings.TrimSpace(c.Param("suggestion")))
	ret.username = strings.TrimSpace(c.Param("username"))

	startStr := strings.TrimSpace(c.Param("start"))
	if startUnix, err := strconv.ParseInt(startStr, 10, 64); err == nil {
		ret.start = time.Unix(startUnix/1000, 0)
	}

	ret.duration, err = time.ParseDuration(c.Param("duration") + "s")
	if err != nil {
		// Something wierd! No guesses.
		ret.duration = 0
	}

	ret.guesses, err = strconv.ParseUint(c.Param("guesses"), 10, 64)
	if err != nil {
		// Something wierd! No guesses.
		ret.guesses = 0
	}

	return ret
}
