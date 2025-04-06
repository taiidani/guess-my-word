package app

import (
	"fmt"
	"guess_my_word/internal/sessions"
	"log/slog"
	"net/http"
)

type indexBag struct {
	baseBag
	Guesser guessBag
	Mode    string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := indexBag{}
	data.Page = "home"

	data.Session = sessions.New(w, r)
	defer func() {
		if err := data.Session.Save(); err != nil {
			slog.Warn("Unable to save session", "error", err)
		}
	}()

	// Assign the mode
	switch {
	// Allow mode selection through the path
	case r.PathValue("mode") != "":
		data.Mode = r.PathValue("mode")
	default:
		data.Mode = "default"
	}
	data.Session.Mode = data.Mode

	// Load list data for the current mode
	// This also validates that it is an existing mode
	var err error
	data.List, err = listStore.GetList(r.Context(), data.Mode)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, fmt.Errorf("could not load list %q: %s", data.Mode, err))
		return
	}

	// Generate the word for the day
	tm := data.Session.DateUser()
	word, err := wordStore.GetForDay(r.Context(), tm, data.Session.Mode)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	data.Guesser = fillGuessBag(data.Session.Current(), wordStore, word)

	renderHtml(w, http.StatusOK, "index.gohtml", data)
}
