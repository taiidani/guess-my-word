package actions

import (
	"guess_my_word/internal/words"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type hint struct {
	After  string    `form:"after"`
	Before string    `form:"before"`
	Mode   string    `form:"mode"`
	Start  time.Time `form:"start" time_format:"unix"`
	TZ     int       `form:"tz"`
}

type hintReply struct {
	Word  string `json:"word"`
	Error string `json:"error"`
}

const (
	// ErrInvalidRequest is emitted when the request payload is malformed
	ErrInvalidRequest = "Invalid request format received"

	// ErrEmptyBeforeAfter is emitted when the before/after have not been provided
	ErrEmptyBeforeAfter = "You need to at least guess the before and after first!"
)

// HintHandler is an API handler to provide a hint to a user.
func HintHandler(c *gin.Context) {
	hint := hint{}
	reply := hintReply{}

	// Validate the guess
	if err := c.ShouldBind(&hint); err != nil {
		log.Println("Invalid request received: ", err)
		reply.Error = ErrInvalidRequest
	} else if len(strings.TrimSpace(hint.Before)) == 0 || len(strings.TrimSpace(hint.After)) == 0 {
		reply.Error = ErrEmptyBeforeAfter
	} else if hint.Start.Unix() == 0 {
		reply.Error = ErrInvalidStartTime
	}
	if reply.Error != "" {
		c.JSON(200, reply)
		return
	}

	// Generate the word for the day
	word, err := words.GetForDay(c, convertUTCToUser(hint.Start, hint.TZ), hint.Mode)
	if err != nil {
		reply.Error = err.Error()
		c.JSON(500, reply)
		return
	}

	reply.Word = getWordHint(hint, word)

	c.JSON(200, reply)
}

func getWordHint(h hint, word string) string {
	minWord := minof(len(h.After), len(h.Before))
	for i := 0; i < minWord; i++ {
		if h.After[i] != h.Before[i] {
			return word[0 : i+1]
		}
	}

	return word[0 : minWord+1]
}

func minof(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
