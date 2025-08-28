package app

import (
	"context"
	"embed"
	"guess_my_word/internal/model"
	"guess_my_word/internal/sessions"
	"net/http"
	"os"
	"time"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	dev       = os.Getenv("DEV") == "true"
)

func SetupStores(l listClient, w wordClient) {
	listStore = l
	wordStore = w
}

//go:embed templates
var templates embed.FS

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r chi.Router) error {
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	r.Use(middleware.Logger)
	r.Use(standardHeadersMiddleware)
	r.Use(sentryHandler.Handle)

	r.Get("/", IndexHandler)
	r.Get("/assets/*", assetsHandler)
	r.Get("/mode/{mode}", IndexHandler)
	r.Get("/about", AboutHandler)
	r.Get("/ping", PingHandler)
	r.Get("/stats", StatsHandler)
	r.Post("/guess", GuessHandler)
	r.Get("/hint", HintHandler)
	r.Post("/reset", ResetHandler)

	r.Route("/api", func(g chi.Router) {
		g.Get("/lists", ListsHandler)
		g.Get("/list", ListHandler)
		g.Post("/list", ListHandler)
		g.Put("/list", ListHandler)
		g.Delete("/list", ListHandler)
		g.Get("/seed", SeedHandler)
	})

	return nil
}

var fnPopulateSessionData func(s *sessions.Session) = func(s *sessions.Session) {}

func startSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	ret := sessions.New(w, r)
	fnPopulateSessionData(ret)

	return ret, nil
}
