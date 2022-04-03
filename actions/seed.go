package actions

import (
	"context"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type seed struct {
	Mode string `form:"mode"`
	TZ   int    `form:"tz"`
}

type seedReply struct {
	Error string `json:"error,omitempty"`
}

// SeedHandler is an internal API handler for pre-populating data to test with.
func SeedHandler(c *gin.Context) {
	seed := seed{}
	reply := seedReply{}

	// Validate the form
	if err := c.ShouldBind(&seed); err != nil {
		log.Println("Invalid request received: ", err)
		reply.Error = ErrInvalidRequest
	}

	if seed.Mode == "" {
		seed.Mode = "default"
	}
	if seed.TZ == 0 {
		seed.TZ = 480 // Pacific
	}

	// Today
	tm := time.Now()
	if err := seedHandlerReply(c, tm, seed.TZ, seed.Mode); err != nil {
		reply.Error = err.Error()
		c.JSON(500, reply)
		return
	}

	// Yesterday
	tm = tm.AddDate(0, 0, -1)
	if err := seedHandlerReply(c, tm, seed.TZ, seed.Mode); err != nil {
		reply.Error = err.Error()
		c.JSON(500, reply)
		return
	}

	c.JSON(200, reply)
}

func seedHandlerReply(ctx context.Context, dt time.Time, tz int, mode string) error {
	// Generate the word for the day
	tm := convertUTCToUser(dt, tz)
	word, err := wordStore.GetForDay(ctx, tm, mode)
	if err != nil {
		return err
	}

	numGuesses := rand.Intn(1000)
	for i := 0; i < numGuesses; i++ {
		guess := model.Guess{Count: rand.Intn(30) + 1}
		word.Guesses = append(word.Guesses, guess)
	}

	return wordStore.SetWord(ctx, datastore.WordKey(mode, tm), word)
}
