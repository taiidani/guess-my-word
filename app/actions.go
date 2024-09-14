package app

import (
	"context"
	"embed"
	"encoding/json"
	"guess_my_word/internal/model"
	"guess_my_word/internal/sessions"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

//go:embed assets
var assets embed.FS

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Serving file", "path", r.URL.Path)
	if dev {
		http.ServeFile(w, r, filepath.Join("app", r.URL.Path))
	} else {
		http.ServeFileFS(w, r, assets, r.URL.Path)
	}
}

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r chi.Router) error {
	r.Use(middleware.Logger)
	r.Use(standardHeadersMiddleware)

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

func renderHtml(w http.ResponseWriter, code int, file string, data any) {
	log := slog.With("name", file, "code", code)

	var t *template.Template
	var err error
	if dev {
		t, err = template.ParseGlob("app/templates/**")
	} else {
		t, err = template.ParseFS(templates, "templates/**")
	}
	if err != nil {
		log.Error("Could not parse templates", "error", err)
		return
	}

	log.Debug("Rendering file", "dev", dev)
	w.WriteHeader(code)
	err = t.ExecuteTemplate(w, file, data)
	if err != nil {
		log.Error("Could not render template", "error", err)
	}
}

func renderJson(w http.ResponseWriter, code int, data any) {
	log := slog.With("code", code)

	log.Debug("Rendering json", "dev", dev)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Error("Could not render template", "error", err)
	}
}

func renderRedirect(w http.ResponseWriter, code int, location string) {
	log := slog.With("code", code)

	log.Debug("Rendering redirect", "dev", dev)
	w.Header().Add("Location", location)
	w.WriteHeader(code)
}
