package app

import (
	"context"
	"embed"
	"guess_my_word/internal/model"
	"guess_my_word/internal/sessions"
	"html/template"
	"io/fs"
	"net/http"
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

type baseBag struct {
	Session *sessions.Session
	Page    string
	List    model.List
}

var (
	listStore listClient
	wordStore wordClient
)

func SetupStores(l listClient, w wordClient) {
	listStore = l
	wordStore = w
}

//go:embed templates
var templates embed.FS

// SetupTemplates will load the HTML templates into gin.
func SetupTemplates(r *gin.Engine) error {
	t, err := template.ParseFS(templates, "templates/**")
	if err != nil {
		return err
	}
	r.SetHTMLTemplate(t)

	return nil
}

//go:embed assets
var assets embed.FS

// SetupAssets will load the static assets into gin.
func SetupAssets(r *gin.Engine) error {
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return err
	}
	r.StaticFS("/assets", http.FS(sub))

	return nil
}

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r *gin.Engine) error {
	r.Use(middlewareStandardHeaders())
	r.GET("/", IndexHandler)
	r.GET("/mode/:mode", IndexHandler)
	r.GET("/about", AboutHandler)
	r.GET("/ping", PingHandler)
	r.GET("/stats", StatsHandler)
	r.POST("/guess", GuessHandler)
	r.GET("/hint", HintHandler)
	r.POST("/reset", ResetHandler)

	g := r.Group("/api")
	g.GET("/lists", ListsHandler)
	g.GET("/list", ListHandler)
	g.POST("/list", ListHandler)
	g.PUT("/list", ListHandler)
	g.DELETE("/list", ListHandler)
	g.GET("/seed", SeedHandler)

	return nil
}

var fnPopulateSessionData func(s *sessions.Session) = func(s *sessions.Session) {}

func startSession(c *gin.Context) (*sessions.Session, error) {
	ret := sessions.New(c)
	fnPopulateSessionData(ret)

	return ret, nil
}
