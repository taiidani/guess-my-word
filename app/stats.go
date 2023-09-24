package app

import (
	_ "embed"
	"guess_my_word/internal/model"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrRevealToday is emitted when the reveal request is for a current or future word
const ErrRevealToday = "It's too early to reveal this word. Please try again later!"

type statsBag struct {
	baseBag
	Yesterday     replyData
	YesterdayHard replyData
	Today         replyData
	TodayHard     replyData
}

// StatsHandler is an HTML handler for pre-populating data to test with.
func StatsHandler(c *gin.Context) {
	// Set up the days
	dateToday := time.Now()
	dateYesterday := dateToday.Add(time.Hour * -24)

	// Generate the word for the day
	wordYesterday, err := wordStore.GetForDay(c, dateYesterday, "default")
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}
	wordYesterdayHard, err := wordStore.GetForDay(c, dateYesterday, "hard")
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}
	wordToday, err := wordStore.GetForDay(c, dateToday, "default")
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}
	wordTodayHard, err := wordStore.GetForDay(c, dateToday, "hard")
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	data := statsBag{}
	data.Page = "stats"
	data.Yesterday = analyzeDay(wordYesterday)
	data.YesterdayHard = analyzeDay(wordYesterdayHard)
	data.Today = analyzeDay(wordToday)
	data.Today.Word = ""
	data.TodayHard = analyzeDay(wordTodayHard)
	data.TodayHard.Word = ""
	c.HTML(http.StatusOK, "stats.gohtml", data)
}

// YesterdayHandler is an HTML handler for pre-populating data to test with.
func YesterdayHandler(c *gin.Context) {
	request, err := parseBodyData(c)
	if err != nil {
		slog.Warn("Unable to parse body data", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	// Subtract one day for yesterday
	dateUser := request.Session.DateUser().Add(time.Hour * -24)

	// Is it too early to reveal the word?
	y, m, d := time.Now().Date()
	cmp := time.Date(y, m, d, 0, 0, 0, 0, dateUser.Location())

	if dateUser.After(cmp) {
		slog.Warn("Too early to reveal word", "date", dateUser)
		c.HTML(http.StatusBadRequest, "error.gohtml", ErrRevealToday)
		return
	}

	// Generate the word for the day
	word, err := wordStore.GetForDay(c, dateUser, request.Session.Mode)
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	data := analyzeDay(word)
	c.HTML(http.StatusOK, "stats.gohtml", data)
}

// TodayHandler is an HTML handler for pre-populating data to test with.
func TodayHandler(c *gin.Context) {
	request, err := parseBodyData(c)
	if err != nil {
		slog.Warn("Unable to parse body data", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	// Generate the word for the day
	word, err := wordStore.GetForDay(c, request.Session.DateUser(), request.Session.Mode)
	if err != nil {
		slog.Warn("Unable to get day", "error", err)
		c.HTML(http.StatusBadRequest, "error.gohtml", err)
		return
	}

	data := analyzeDay(word)

	// Wipe the word from the data, as it's today
	data.Word = ""

	c.HTML(http.StatusOK, "stats.gohtml", data)
}

type replyData struct {
	Word        string
	Completions int
	BestRun     int
	AvgRun      int
}

func analyzeDay(word model.Word) replyData {
	// If no one guessed that day
	if len(word.Guesses) == 0 {
		return replyData{Word: word.Value}
	}

	ret := replyData{
		Word:        word.Value,
		Completions: len(word.Guesses),
		BestRun:     999,
	}

	var guessCount = 0
	for _, item := range word.Guesses {
		guessCount += item.Count

		if item.Count < ret.BestRun {
			ret.BestRun = item.Count
		}
	}

	ret.AvgRun = guessCount / len(word.Guesses)
	return ret
}
