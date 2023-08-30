package sessions

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Session struct {
	// Mode tracks the current word list being guessed
	Mode string `json:"mode"`

	// History tracks the user's activity across each mode
	History map[string]*SessionMode `json:"history"`

	// session holds the internal session record for later saving
	session sessions.Session
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

func New(c *gin.Context) *Session {
	s := sessions.Default(c)

	jsonSession := s.Get("session")
	session := Session{
		Mode:    "default",
		History: map[string]*SessionMode{},
	}
	if jsonSession != nil {
		if err := json.Unmarshal(jsonSession.([]byte), &session); err != nil {
			slog.Warn("Could not parse history", "error", err)
		}
	}

	session.session = s
	return &session
}

func Configure(r *gin.Engine, client sessions.Store) {
	r.Use(sessions.Sessions("guessmyword", client))
}

func (s *Session) Clear() error {
	s.session.Clear()
	return s.session.Save()
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

	s.session.Set("session", h)
	return s.session.Save()
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

func (m *SessionMode) Stale() bool {
	now := time.Now()
	return m.Start.Month() != now.Month() || m.Start.Day() != now.Day()
}
