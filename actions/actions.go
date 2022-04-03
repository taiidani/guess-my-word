package actions

import (
	"context"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"guess_my_word/internal/words"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type listClient interface {
	GetLists(ctx context.Context) ([]string, error)
	GetList(ctx context.Context, name string) (model.List, error)
	CreateList(ctx context.Context, name string, list model.List) error
	DeleteList(ctx context.Context, name string) error
	UpdateList(ctx context.Context, name string, list model.List) error
}

type wordClient interface {
	Validate(context.Context, string) bool
	GetForDay(context.Context, time.Time, string) (model.Word, error)
	GetWord(ctx context.Context, key string) (model.Word, error)
	SetWord(ctx context.Context, key string, word model.Word) error
}

var (
	listStore listClient
	wordStore wordClient
)

func init() {
	var client datastore.Client
	if addr, ok := os.LookupEnv("REDIS_ADDR"); ok {
		client = datastore.NewRedis(addr)
	} else {
		log.Println("WARNING: No REDIS_ADDR env var set. Falling back upon in-memory store")
		client = datastore.NewMemory()
	}

	listStore = words.NewListStore(client)
	wordStore = words.NewWordStore(client)
}

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r *gin.Engine) (err error) {
	r.Use(middlewareStandardHeaders())
	r.GET("/ping", PingHandler)

	g := r.Group("/api")
	g.GET("/guess", GuessHandler)
	g.GET("/hint", HintHandler)
	g.GET("/lists", ListsHandler)
	g.GET("/list", ListHandler)
	g.POST("/list", ListHandler)
	g.PUT("/list", ListHandler)
	g.DELETE("/list", ListHandler)
	g.GET("/seed", SeedHandler)
	g.GET("/stats", StatsHandler)

	// And the websockets
	g.GET("/ws", wsHandler)
	return nil
}

// convertUTCToLocal will take a given time in UTC and convert it to a given user's timezone
// TZ for PDT (-7:00) is a positive 420, so SUBTRACT that from the unix timestamp
func convertUTCToUser(t time.Time, tz int) time.Time {
	ret := t.In(time.FixedZone("User", tz*-1))
	ret = ret.Add(time.Minute * -1 * time.Duration(tz))
	return ret
}
