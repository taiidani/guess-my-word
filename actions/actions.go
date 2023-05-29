package actions

import (
	"context"
	"guess_my_word/internal/model"
	"guess_my_word/internal/sessions"
	"log"
	"strconv"
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

func SetupStores(l listClient, w wordClient) {
	listStore = l
	wordStore = w
}

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r *gin.Engine) (err error) {
	r.Use(middlewareStandardHeaders())
	r.GET("/", IndexHandler)
	r.GET("/ping", PingHandler)
	r.GET("/stats/yesterday", YesterdayHandler)
	r.GET("/stats/today", TodayHandler)
	r.POST("/mode", ModeSetHandler)
	r.POST("/guess", GuessHandler)
	r.GET("/hint", HintHandler)

	g := r.Group("/api")
	g.GET("/lists", ListsHandler)
	g.GET("/list", ListHandler)
	g.POST("/list", ListHandler)
	g.PUT("/list", ListHandler)
	g.DELETE("/list", ListHandler)

	return nil
}

type bodyData struct {
	TZ      int               // The timezone offset for the user, in milliseconds
	Session *sessions.Session // The session for the user
}

func parseBodyData(c *gin.Context) (bodyData, error) {
	ret := bodyData{}
	ret.Session = sessions.New(c)

	tz, err := strconv.ParseInt(c.Request.URL.Query().Get("tz"), 10, 64)
	if err != nil {
		log.Println("ERROR: Could not parse timezone: ", err)
		tz = 0
	}
	ret.TZ = int(tz)

	return ret, nil
}
