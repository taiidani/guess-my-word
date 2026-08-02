package hangman

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
)

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r chi.Router) error {
	r.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{}))
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
