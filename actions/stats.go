package actions

import (
	"context"
	"guess_my_word/internal/model"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type stats struct {
	Date     time.Time `form:"date" time_format:"unix"` // The unix date, no timestamp
	dateUser time.Time // The date under the user's timezone
	TZ       int       `form:"tz"`
	Mode     string    `form:"mode"`
}

type statsReply struct {
	Today model.Word `json:"today"`
	Word  model.Word `json:"word"`
	Error string     `json:"error,omitempty"`
}

// ErrRevealToday is emitted when the reveal request is for a current or future word
const ErrRevealToday = "It's too early to reveal this word. Please try again later!"

// StatsHandler is an internal API handler for pre-populating data to test with.
func StatsHandler(c *gin.Context) {
	request := stats{}
	reply := statsReply{}

	// Validate the form
	if err := c.ShouldBind(&request); err != nil {
		log.Println("Invalid request received: ", err)
		reply.Error = ErrInvalidRequest
	}

	if request.Date.Unix() == 0 {
		reply.Error = ErrInvalidStartTime
	} else {
		log.Println("TZ:", request.TZ)
		request.dateUser = convertUTCToUser(request.Date, request.TZ)
		log.Printf("Requested date is: %s", request.dateUser)

		y, m, d := time.Now().Date()
		cmp := time.Date(y, m, d, 0, 0, 0, 0, request.dateUser.Location())

		if request.dateUser.After(cmp) {
			log.Printf("Compared date was: %s", request.dateUser)
			reply.Error = ErrRevealToday
		}
	}

	if reply.Error != "" {
		c.JSON(400, reply)
		return
	}

	reply = refreshStats(c, request)
	c.JSON(200, reply)
}

func refreshStats(ctx context.Context, request stats) statsReply {
	reply := statsReply{}

	// Generate the word for the day
	word, err := wordStore.GetForDay(ctx, request.dateUser, request.Mode)
	if err != nil {
		reply.Error = err.Error()
		return reply
	}

	reply.Word = word

	// Now for today's stats. Similar, but without the word information!
	todayTm := request.dateUser.AddDate(0, 0, 1)
	word, err = wordStore.GetForDay(ctx, todayTm, request.Mode)
	if err != nil {
		reply.Error = err.Error()
		return reply
	}

	reply.Today = word
	reply.Today.Value = ""
	return reply
}
