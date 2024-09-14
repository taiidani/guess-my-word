package app

import (
	"guess_my_word/internal/sessions"
	"log/slog"
	"net/http"
)

const (
	// ErrInvalidRequest is emitted when the request payload is malformed
	ErrInvalidRequest = "Invalid request format received"

	// ErrEmptyBeforeAfter is emitted when the before/after have not been provided
	ErrEmptyBeforeAfter = "You need to at least guess the before and after first!"
)

// HintHandler is an API handler to provide a hint to a user.
func HintHandler(w http.ResponseWriter, r *http.Request) {
	session, err := startSession(w, r)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	// Generate the word for the day
	h := session.Current()
	word, err := wordStore.GetForDay(r.Context(), h.DateUser(), session.Mode)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	hintWord := getWordHint(h, word.Value)
	renderHtml(w, http.StatusOK, "raw.gohtml", "The word starts with: "+hintWord)
}

func getWordHint(h *sessions.SessionMode, word string) string {
	letters := h.CommonGuessPrefix()

	// Don't return the whole word if there's only one letter left!
	if len(letters) >= len(word)-1 {
		return word[0 : len(word)-1]
	}

	return word[0 : len(letters)+1]
}
