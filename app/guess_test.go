package app

import (
	"errors"
	"guess_my_word/internal/model"
	"guess_my_word/internal/sessions"
	"guess_my_word/internal/words"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_GuessHandler(t *testing.T) {
	router := setupRouter()
	renderer := template.Must(template.ParseFS(templates, "templates/**"))

	// Force the end time for correct guesses to be predictable
	mockEndTime := time.Now()
	fnSetEndTime = func() *time.Time { return &mockEndTime }

	render := func(name string, data any) string {
		buf := strings.Builder{}
		_ = renderer.ExecuteTemplate(&buf, name, data)
		return buf.String()
	}

	tests := []struct {
		name    string
		session func(s *sessions.Session)
		post    url.Values
		want    string
		wantErr bool
	}{
		{
			name: "After",
			post: url.Values{"word": {"power"}},
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("guesser.gohtml", guessBag{
				Session: &sessions.SessionMode{
					Start:  time.Unix(1587930259, 0), // generates "belong"
					Before: []string{},
					After:  []string{"power"},
				},
				ProgressBefore:        0,
				ProgressRange:         20,
				ProgressAfter:         80,
				ProgressAfterAbsolute: 20,
			}),
		},
		{
			name: "Before",
			post: url.Values{"word": {"apple"}},
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("guesser.gohtml", guessBag{
				Session: &sessions.SessionMode{
					Start:  time.Unix(1587930259, 0), // generates "belong"
					Before: []string{"apple"},
					After:  []string{},
				},
				ProgressBefore:        0,
				ProgressRange:         100,
				ProgressAfter:         0,
				ProgressAfterAbsolute: 100,
			}),
		},
		{
			name: "Correct",
			post: url.Values{"word": {"belong"}},
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{"apple"},
						After:  []string{},
					},
				}
			},
			want: render("guesser.gohtml", guessBag{
				Session: &sessions.SessionMode{
					Start:  time.Unix(1587930259, 0), // generates "belong"
					End:    &mockEndTime,
					Before: []string{"apple"},
					After:  []string{},
					Answer: "belong",
				},
				Stats: model.WordStats{
					BestRun: 2,
				},
			}),
		},
		{
			name: "Invalid word",
			post: url.Values{"word": {"asdf"}},
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("error.gohtml", errorBag{Message: errors.New(ErrInvalidWord)}),
		},
		{
			name: "Empty word",
			post: url.Values{"word": {" "}},
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("error.gohtml", errorBag{Message: errors.New(ErrEmptyGuess)}),
		},
		{
			name: "Correct Tomorrow",
			post: url.Values{"word": {"roll"}},
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1588030259, 0), // generates "roll"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("guesser.gohtml", guessBag{
				Session: &sessions.SessionMode{
					Start:  time.Unix(1588030259, 0), // generates "roll"
					End:    &mockEndTime,
					Before: []string{},
					After:  []string{},
					Answer: "roll",
				},
				Stats: model.WordStats{
					BestRun: 1,
				},
			}),
		},
		{
			name: "Correct Yesterday",
			post: url.Values{"word": {"laundry"}},
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587830259, 0), // generates "laundry"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("guesser.gohtml", guessBag{
				Session: &sessions.SessionMode{
					Start:  time.Unix(1587830259, 0), // generates "laundry"
					End:    &mockEndTime,
					Before: []string{},
					After:  []string{},
					Answer: "laundry",
				},
				Stats: model.WordStats{
					BestRun: 1,
				},
			}),
		},
		{
			name: "Correct Hard",
			post: url.Values{"word": {"teth"}},
			session: func(s *sessions.Session) {
				s.Mode = "hard"
				s.History = map[string]*sessions.SessionMode{
					"hard": {
						Start:  time.Unix(1587930259, 0), // generates "teth"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("guesser.gohtml", guessBag{
				Session: &sessions.SessionMode{
					Start:  time.Unix(1587930259, 0), // generates "teth"
					End:    &mockEndTime,
					Before: []string{},
					After:  []string{},
					Answer: "teth",
				},
				Stats: model.WordStats{
					BestRun: 1,
				},
			}),
		},
		{
			name: "Correct Hard Yesterday",
			post: url.Values{"word": {"tayra"}},
			session: func(s *sessions.Session) {
				s.Mode = "hard"
				s.History = map[string]*sessions.SessionMode{
					"hard": {
						Start:  time.Unix(1587830259, 0), // generates "tayra"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("guesser.gohtml", guessBag{
				Session: &sessions.SessionMode{
					Start:  time.Unix(1587830259, 0), // generates "tayra"
					End:    &mockEndTime,
					Before: []string{},
					After:  []string{},
					Answer: "tayra",
				},
				Stats: model.WordStats{
					BestRun: 1,
				},
			}),
		},
		{
			name: "Invalid request",
			post: url.Values{},
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{},
						After:  []string{},
					},
				}
			},
			want: render("error.gohtml", errorBag{Message: errors.New(ErrEmptyGuess)}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fnPopulateSessionData = tt.session

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/guess", nil)
			req.PostForm = tt.post
			router.ServeHTTP(w, req)

			got := w.Body.String()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guess() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_fillGuessBag(t *testing.T) {
	type args struct {
		s    *sessions.SessionMode
		word model.Word
	}
	tests := []struct {
		name string
		args args
		want guessBag
	}{
		{
			name: "empty",
			args: args{
				s: &sessions.SessionMode{
					BeforeI: 0,
					AfterI:  words.ScrabbleDictionary.Size(), // Dictionary size
				},
			},
			want: guessBag{
				Session: &sessions.SessionMode{
					BeforeI: 0,
					AfterI:  words.ScrabbleDictionary.Size(), // Dictionary size
				},
				ProgressBefore:        0,
				ProgressRange:         100,
				ProgressAfter:         0,
				ProgressAfterAbsolute: 100,
			},
		},
		{
			name: "before populated",
			args: args{
				s: &sessions.SessionMode{
					BeforeI: 20000,
					AfterI:  words.ScrabbleDictionary.Size(), // Dictionary size
				},
			},
			want: guessBag{
				Session: &sessions.SessionMode{
					BeforeI: 0,
					AfterI:  words.ScrabbleDictionary.Size(), // Dictionary size
				},
				ProgressBefore:        7,
				ProgressRange:         93,
				ProgressAfter:         0,
				ProgressAfterAbsolute: 100,
			},
		},
		{
			name: "after populated",
			args: args{
				s: &sessions.SessionMode{
					BeforeI: 0,
					AfterI:  words.ScrabbleDictionary.Size(), // Dictionary size
				},
			},
			want: guessBag{
				Session: &sessions.SessionMode{
					BeforeI: 0,
					AfterI:  words.ScrabbleDictionary.Size(), // Dictionary size
				},
				ProgressBefore:        0,
				ProgressRange:         100,
				ProgressAfter:         0,
				ProgressAfterAbsolute: 100,
			},
		},
		{
			name: "both populated",
			args: args{
				s: &sessions.SessionMode{
					BeforeI: 0,
					AfterI:  words.ScrabbleDictionary.Size(), // Dictionary size
				},
			},
			want: guessBag{
				Session: &sessions.SessionMode{
					BeforeI: 0,
					AfterI:  words.ScrabbleDictionary.Size(), // Dictionary size
				},
				ProgressBefore:        0,
				ProgressRange:         100,
				ProgressAfter:         0,
				ProgressAfterAbsolute: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillGuessBag(tt.args.s, tt.args.word); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fillGuessBag() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
