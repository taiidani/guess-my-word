package app

import (
	_ "embed"
	"errors"
	"guess_my_word/internal/model"
	"log/slog"
	"net/http"
	"time"
)

// ErrRevealToday is emitted when the reveal request is for a current or future word
const ErrRevealToday = "It's too early to reveal this word. Please try again later!"

type statsBag struct {
	baseBag
	Yesterday     model.WordStats
	YesterdayHard model.WordStats
	Today         model.WordStats
	TodayHard     model.WordStats
}

// StatsHandler is an HTML handler for pre-populating data to test with.
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := startSession(w, r)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	// Set up the days
	dateToday := session.DateUser()
	dateYesterday := dateToday.Add(time.Hour * -24)

	// Generate the word for the day
	wordYesterday, err := wordStore.GetForDay(r.Context(), dateYesterday, "default")
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	wordYesterdayHard, err := wordStore.GetForDay(r.Context(), dateYesterday, "hard")
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	wordToday, err := wordStore.GetForDay(r.Context(), dateToday, "default")
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}
	wordTodayHard, err := wordStore.GetForDay(r.Context(), dateToday, "hard")
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	data := statsBag{}
	data.Session = session
	data.Page = "stats"
	data.Yesterday = wordYesterday.Stats()
	data.YesterdayHard = wordYesterdayHard.Stats()
	data.Today = wordToday.Stats()
	data.Today.Word = ""
	data.TodayHard = wordTodayHard.Stats()
	data.TodayHard.Word = ""
	renderHtml(w, http.StatusOK, "stats.gohtml", data)
}

// YesterdayHandler is an HTML handler for pre-populating data to test with.
func YesterdayHandler(w http.ResponseWriter, r *http.Request) {
	session, err := startSession(w, r)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	// Subtract one day for yesterday
	dateUser := session.DateUser().Add(time.Hour * -24)

	// Is it too early to reveal the word?
	y, m, d := time.Now().Date()
	cmp := time.Date(y, m, d, 0, 0, 0, 0, dateUser.Location())

	if dateUser.After(cmp) {
		slog.Warn("Too early to reveal word", "date", dateUser)
		errorResponse(w, http.StatusBadRequest, errors.New(ErrRevealToday))
		return
	}

	// Generate the word for the day
	word, err := wordStore.GetForDay(r.Context(), dateUser, session.Mode)
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	data := word.Stats()
	renderHtml(w, http.StatusOK, "stats.gohtml", data)
}

// TodayHandler is an HTML handler for pre-populating data to test with.
func TodayHandler(w http.ResponseWriter, r *http.Request) {
	session, err := startSession(w, r)
	if err != nil {
		slog.Warn("Unable to start session", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	// Generate the word for the day
	word, err := wordStore.GetForDay(r.Context(), session.DateUser(), session.Mode)
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		errorResponse(w, http.StatusBadRequest, err)
		return
	}

	data := word.Stats()

	// Wipe the word from the data, as it's today
	data.Word = ""

	renderHtml(w, http.StatusOK, "stats.gohtml", data)
}
