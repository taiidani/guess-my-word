package app

import (
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/sessions"
	"guess_my_word/internal/words"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/quasoft/memstore"
)

func init() {
	client := datastore.NewMemory()
	listStore = words.NewListStore(client)
	wordStore = words.NewWordStore(client)
}

func setupRouter(t *testing.T) chi.Router {
	t.Setenv("DEV", "true")

	r := chi.NewRouter()

	sessionClient := memstore.NewMemStore([]byte("secret"))
	sessions.Configure(r, "guessmyword", sessionClient)

	_ = AddHandlers(r)
	return r
}
