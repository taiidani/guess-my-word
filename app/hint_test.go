package app

import (
	"guess_my_word/internal/sessions"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_HintHandler(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name    string
		session func(s *sessions.Session)
		want    string
		wantErr bool
	}{
		{
			name: "Not even close",
			session: func(s *sessions.Session) {
				s.Mode = "hard"
				s.History = map[string]*sessions.SessionMode{
					"hard": {
						Start:  time.Unix(1587930259, 0), // generates "teth"
						Before: []string{"apple"},
						After:  []string{"zoo"},
					},
				}
			},
			want: `The word starts with: t
`,
		},
		{
			name: "1 character",
			session: func(s *sessions.Session) {
				s.Mode = "hard"
				s.History = map[string]*sessions.SessionMode{
					"hard": {
						Start:  time.Unix(1587930259, 0), // generates "teth"
						Before: []string{"train"},
						After:  []string{"zoo"},
					},
				}
			},
			want: `The word starts with: t
`,
		},
		{
			name: "2 characters",
			session: func(s *sessions.Session) {
				s.Mode = "hard"
				s.History = map[string]*sessions.SessionMode{
					"hard": {
						Start:  time.Unix(1587930259, 0), // generates "teth"
						Before: []string{"tear"},
						After:  []string{"train"},
					},
				}
			},
			want: `The word starts with: te
`,
		},
		{
			name: "Almost there",
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{"belo"},
						After:  []string{"belonging"},
					},
				}
			},
			want: `The word starts with: belon
`,
		},
		{
			name: "1 letter left",
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{"belon"},
						After:  []string{"belony"},
					},
				}
			},
			want: `The word starts with: belon
`,
		},
		{
			name: "Empty before",
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{""},
						After:  []string{"belonging"},
					},
				}
			},
			want: `The word starts with: b
`,
		},
		{
			name: "Empty after",
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{"belo"},
						After:  []string{""},
					},
				}
			},
			want: `The word starts with: b
`,
		},
		{
			name: "Empty both",
			session: func(s *sessions.Session) {
				s.Mode = "default"
				s.History = map[string]*sessions.SessionMode{
					"default": {
						Start:  time.Unix(1587930259, 0), // generates "belong"
						Before: []string{""},
						After:  []string{""},
					},
				}
			},
			want: `The word starts with: b
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fnPopulateTestSessionData = tt.session

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/hint", nil)
			router.ServeHTTP(w, req)

			got := w.Body.String()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guess() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
