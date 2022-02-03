package actions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func Test_HintHandler(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name    string
		request url.Values
		want    hintReply
		wantErr bool
	}{
		{
			name: "Not even close",
			request: url.Values{
				"after":  {"apple"},
				"before": {"zoo"},
				"start":  {"1587930259"}, // generates "teth"
				"mode":   {"hard"},
			},
			want: hintReply{
				Word:  "t",
				Error: "",
			},
		},
		{
			name: "1 character",
			request: url.Values{
				"after":  {"zoo"},
				"before": {"train"},
				"start":  {"1587930259"}, // generates "teth"
				"mode":   {"hard"},
			},
			want: hintReply{
				Word:  "t",
				Error: "",
			},
		},
		{
			name: "2 characters",
			request: url.Values{
				"after":  {"tear"},
				"before": {"train"},
				"start":  {"1587930259"}, // generates "teth"
				"mode":   {"hard"},
			},
			want: hintReply{
				Word:  "te",
				Error: "",
			},
		},
		{
			name: "Almost there",
			request: url.Values{
				"after":  {"belo"},
				"before": {"belonging"},
				"start":  {"1587930259"}, // generates "belong"
				"mode":   {"default"},
			},
			want: hintReply{
				Word:  "belon",
				Error: "",
			},
		},
		{
			name: "Empty before",
			request: url.Values{
				"after":  {""},
				"before": {"belonging"},
				"start":  {"1587930259"},
				"mode":   {"default"},
			},
			want: hintReply{
				Word:  "",
				Error: ErrEmptyBeforeAfter,
			},
		},
		{
			name: "Empty after",
			request: url.Values{
				"after":  {"belo"},
				"before": {""},
				"start":  {"1587930259"},
				"mode":   {"default"},
			},
			want: hintReply{
				Word:  "",
				Error: ErrEmptyBeforeAfter,
			},
		},
		{
			name: "Empty both",
			request: url.Values{
				"after":  {""},
				"before": {""},
				"start":  {"1587930259"},
				"mode":   {"default"},
			},
			want: hintReply{
				Word:  "",
				Error: ErrEmptyBeforeAfter,
			},
		},
		{
			name: "Invalid time",
			request: url.Values{
				"after":  {"belo"},
				"before": {"belonging"},
				"start":  {"0"}, // OMG WUT
				"mode":   {"default"},
			},
			want: hintReply{
				Word:  "",
				Error: ErrInvalidStartTime,
			},
		},
		{
			name: "Invalid request",
			request: url.Values{
				"start": {"bar"},
			},
			want: hintReply{
				Word:  "",
				Error: ErrInvalidRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/hint?"+tt.request.Encode(), nil)
			router.ServeHTTP(w, req)

			got := hintReply{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("guess() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
