package actions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
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

func Test_generateWord(t *testing.T) {
	tests := []struct {
		name    string
		seed    time.Time
		want    string
		wantErr bool
	}{
		{
			name: "Date yesterday",
			seed: time.Date(2020, time.January, 26, 2, 0, 0, 0, time.UTC),
			want: "power",
		},
		{
			name: "Date tweak yesterday",
			seed: time.Date(2020, time.January, 27, 2, 0, 0, 0, time.UTC).UTC().AddDate(0, 0, -1),
			want: "power",
		},
		{
			name: "Unix yesterday",
			seed: time.Unix(1580083199, 0), // Sun Jan 26 23:59:59 2020 UTC
			want: "power",
		},
		{
			name: "Date today",
			seed: time.Date(2020, time.January, 27, 2, 0, 0, 0, time.UTC),
			want: "tell",
		},
		{
			name: "Date tweak today",
			seed: time.Date(2020, time.January, 26, 2, 0, 0, 0, time.UTC).UTC().AddDate(0, 0, 1),
			want: "tell",
		},
		{
			name: "Unix today",
			seed: time.Unix(1580083201, 0), // Mon Jan 27 00:00:01 2020 UTC
			want: "tell",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateWord(tt.seed, words)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateWord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
