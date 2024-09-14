package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

type Session struct {
	// Mode tracks the current word list being guessed
	Mode string `json:"mode"`

	// History tracks the user's activity across each mode
	History map[string]*SessionMode `json:"history"`

	w http.ResponseWriter

	r *http.Request

	session *sessions.Session
}

type SessionMode struct {
	// Start tracks when the user began guessing
	Start time.Time `json:"start"`

	// End tracks when the correct guess was made
	End *time.Time `json:"end"`

	// Before tracks all guesses that the correct word is before
	Before []string `json:"before"`

	// After tracks all guesses that the correct word is after
	After []string `json:"after"`

	// Answer stores the correct answer once it has been found
	Answer string `json:"answer"`
}

func New(w http.ResponseWriter, r *http.Request) *Session {
	sessAny := r.Context().Value("session")
	s := sessAny.(*sessions.Session)

	session := Session{
		Mode:    "default",
		History: map[string]*SessionMode{},
		r:       r,
		w:       w,
	}
	if jsonSession, ok := s.Values["session"]; ok {
		if err := json.Unmarshal(jsonSession.([]byte), &session); err != nil {
			slog.Warn("Could not parse history", "error", err)
		}
	}

	session.session = s
	return &session
}

func Configure(r chi.Router, name string, client sessions.Store) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := client.Get(r, "guessmyword")
			if err != nil {
				slog.Warn("Unable to load session", "error", err)
			}

			ctx := context.WithValue(r.Context(), "session", sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func (s *Session) Clear() error {
	for key := range s.session.Values {
		delete(s.session.Values, key)
	}

	return s.session.Save(s.r, s.w)
}

func (s *Session) Current() *SessionMode {
	if _, ok := s.History[s.Mode]; !ok {
		s.History[s.Mode] = &SessionMode{
			Start:  time.Now(),
			Before: []string{},
			After:  []string{},
		}
	}

	return s.History[s.Mode]
}

func (s *Session) DateUser() time.Time {
	m := s.Current()
	return m.Start
}

func (s *Session) Save() error {
	h, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("could not serialize to session: %s", err)
	}

	s.session.Values["session"] = h
	return s.session.Save(s.r, s.w)
}

func (m *SessionMode) GuessCount() int {
	count := len(m.Before) + len(m.After)
	if m.Answer != "" {
		count++
	}

	return count
}

// CommonGuessPrefix will find the shared letters between the closest Before
// and After guess.
//
// For example, "trace" and "tree" will produce "tr".
//
// An empty string is returned if there are no common letters.
func (m *SessionMode) CommonGuessPrefix() string {
	// Not enough guesses for a common prefix
	if len(m.After) == 0 || len(m.Before) == 0 {
		return ""
	}

	before := m.Before[len(m.Before)-1]
	after := m.After[0]

	minWord := min(len(after), len(before))
	for i := 0; i < minWord; i++ {
		if after[i] != before[i] {
			return before[0:i]
		}
	}

	return before[0:minWord]
}

func (m *SessionMode) DateUser() time.Time {
	return m.Start
}

var remainingSeedTime time.Time

func (m *SessionMode) RemainingTime() string {
	now := time.Now()
	if !remainingSeedTime.IsZero() {
		now = remainingSeedTime
	}

	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)

	// The time between tomorrow (midnight) and right now
	interval := tomorrow.Sub(now)

	// Number of hours, floating point precision
	hours := math.Floor(interval.Hours())

	// Subtract the truncated hours=>minutes from the total minutes
	minutes := math.Floor(interval.Minutes() - (hours * 60))

	// Truncate the hours & minutes, and print
	hoursStr := "hours"
	if hours >= 1 && hours < 2 {
		hoursStr = "hour"
	}

	minutesStr := "minutes"
	if minutes >= 1 && minutes < 2 {
		minutesStr = "minute"
	}
	return fmt.Sprintf("%0.f %s, %0.f %s", hours, hoursStr, minutes, minutesStr)
}

func (m *SessionMode) Stale() bool {
	now := time.Now()
	return m.Start.Month() != now.Month() || m.Start.Day() != now.Day()
}
