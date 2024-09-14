package app

import (
	"guess_my_word/internal/sessions"
	"log/slog"
	"net/http"
)

type aboutBag struct {
	baseBag
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := aboutBag{}
	data.Page = "about"

	data.Session = sessions.New(w, r)
	defer func() {
		if err := data.Session.Save(); err != nil {
			slog.Warn("Unable to save session", "error", err)
		}
	}()

	renderHtml(w, http.StatusOK, "about.gohtml", data)
}
