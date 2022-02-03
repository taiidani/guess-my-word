package actions

import (
	"context"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"guess_my_word/internal/words"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type wordClient interface {
	Validate(context.Context, string) bool
	GetForDay(context.Context, time.Time, string) (model.Word, error)
}

var dataStore words.Store
var wordStore wordClient

func init() {
	dataStore = datastore.NewRedis(os.Getenv("REDIS_ADDR"))
	wordStore = words.NewWordStore(dataStore)
}

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r *gin.Engine) (err error) {
	r.Use(middlewareStandardHeaders())
	r.GET("/api/reveal", RevealHandler)
	r.GET("/ping", PingHandler)
	r.GET("/api/guess", GuessHandler)
	r.GET("/api/hint", HintHandler)
	return nil
}

// convertUTCToLocal will take a given time in UTC and convert it to a given user's timezone
// TZ for PDT (-7:00) is a positive 420, so SUBTRACT that from the unix timestamp
func convertUTCToUser(t time.Time, tz int) time.Time {
	ret := t.In(time.FixedZone("User", tz*-1))
	ret = ret.Add(time.Minute * -1 * time.Duration(tz))
	return ret
}
