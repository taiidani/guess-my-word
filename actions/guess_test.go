package actions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func Test_GuessHandler(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name    string
		request url.Values
		want    guessReply
		wantErr bool
	}{
		{
			name: "Before",
			request: url.Values{
				"word":  {"power"},
				"start": {"1587930259"},
				"mode":  {"default"},
			},
			want: guessReply{
				Guess:   "power",
				Correct: false,
				After:   false,
				Before:  true,
				Error:   "",
			},
		},
		{
			name: "Before",
			request: url.Values{
				"word":  {"apple"},
				"start": {"1587930259"},
				"mode":  {"default"},
			},
			want: guessReply{
				Guess:   "apple",
				Correct: false,
				After:   true,
				Before:  false,
				Error:   "",
			},
		},
		{
			name: "Correct",
			request: url.Values{
				"word":  {"belong"},
				"start": {"1587930259"},
				"mode":  {"default"},
			},
			want: guessReply{
				Guess:   "belong",
				Correct: true,
				After:   false,
				Before:  false,
				Error:   "",
			},
		},
		{
			name: "Invalid word",
			request: url.Values{
				"word":  {"asdf"},
				"start": {"1587930259"},
				"mode":  {"default"},
			},
			want: guessReply{
				Guess:   "asdf",
				Correct: false,
				After:   false,
				Before:  false,
				Error:   ErrInvalidWord,
			},
		},
		{
			name: "Empty word",
			request: url.Values{
				"word":  {" "},
				"start": {"1587930259"},
				"mode":  {"default"},
			},
			want: guessReply{
				Guess:   "",
				Correct: false,
				After:   false,
				Before:  false,
				Error:   ErrEmptyGuess,
			},
		},
		{
			name: "Invalid Time",
			request: url.Values{
				"word":  {"power"},
				"start": {"0"},
				"mode":  {"default"},
			},
			want: guessReply{
				Guess:   "power",
				Correct: false,
				After:   false,
				Before:  false,
				Error:   ErrInvalidStartTime,
			},
		},
		{
			name: "Correct Tomorrow",
			request: url.Values{
				"word":  {"roll"},
				"start": {"1588030259"},
				"mode":  {"default"},
			},
			want: guessReply{
				Guess:   "roll",
				Correct: true,
				After:   false,
				Before:  false,
				Error:   "",
			},
		},
		{
			name: "Correct Yesterday",
			request: url.Values{
				"word":  {"laundry"},
				"start": {"1587830259"},
				"mode":  {"default"},
			},
			want: guessReply{
				Guess:   "laundry",
				Correct: true,
				After:   false,
				Before:  false,
				Error:   "",
			},
		},
		{
			name: "Correct Hard",
			request: url.Values{
				"word":  {"teth"},
				"start": {"1587930259"},
				"mode":  {"hard"},
			},
			want: guessReply{
				Guess:   "teth",
				Correct: true,
				After:   false,
				Before:  false,
				Error:   "",
			},
		},
		{
			name: "Correct Hard Yesterday",
			request: url.Values{
				"word":  {"tayra"},
				"start": {"1587830259"},
				"mode":  {"hard"},
			},
			want: guessReply{
				Guess:   "tayra",
				Correct: true,
				After:   false,
				Before:  false,
				Error:   "",
			},
		},
		{
			name: "Invalid request",
			request: url.Values{
				"start": {"bar"},
			},
			want: guessReply{
				Error: ErrInvalidRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/guess?"+tt.request.Encode(), nil)
			router.ServeHTTP(w, req)

			got := guessReply{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guess() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
