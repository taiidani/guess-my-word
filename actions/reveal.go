package actions

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type reveal struct {
	Date     time.Time `form:"date" time_format:"unix"` // The unix date, no timestamp
	dateUser time.Time // The date under the user's timezone
	TZ       int       `form:"tz"`
	Mode     string    `form:"mode"`
}

type revealReply struct {
	Word  string `json:"word"`
	Error string `json:"error,omitempty"`
}

// ErrRevealToday is emitted when the reveal request is for a current or future word
const ErrRevealToday = "It's too early to reveal this word. Please try again later!"

// RevealHandler reveals the word for a given day
// ...but not today :)
func RevealHandler(c *gin.Context) {
	reveal := reveal{}
	reply := revealReply{}

	// Validate the reveal
	if err := c.ShouldBind(&reveal); err != nil {
		log.Println("Invalid request received: ", err)
		reply.Error = ErrInvalidRequest
	} else if reveal.Date.Unix() == 0 {
		reply.Error = ErrInvalidStartTime
	} else {
		reveal.dateUser = convertUTCToUser(reveal.Date, reveal.TZ)
		log.Printf("Requested date is: %s", reveal.dateUser)

		y, m, d := time.Now().Date()
		cmp := time.Date(y, m, d, 0, 0, 0, 0, reveal.dateUser.Location())

		if reveal.dateUser.After(cmp) {
			log.Printf("Compared date was: %s", reveal.dateUser)
			reply.Error = ErrRevealToday
		}
	}

	if reply.Error != "" {
		c.JSON(200, reply)
		return
	}

	// Generate the word for the day
	word, err := wordStore.GetForDay(c, reveal.dateUser, reveal.Mode)
	if err != nil {
		reply.Error = err.Error()
		c.JSON(500, reply)
		return
	}

	reply.Word = word
	c.JSON(200, reply)
}
