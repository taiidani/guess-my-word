package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"guess_my_word/internal/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestRevealHandler(t *testing.T) {
	type args struct {
		request string
	}
	tests := []struct {
		name          string
		args          args
		mockGetForDay func(ctx context.Context, tm time.Time, mode string) (model.Word, error)
		want          revealReply
		wantCode      int
	}{
		{
			name: "success",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{Value: "theword"}, nil
			},
			args: args{
				request: "/reveal",
			},
			want:     revealReply{Word: model.Word{Value: "theword"}},
			wantCode: 200,
		},
		{
			name: "error-invalid-request",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{Value: "theword"}, nil
			},
			args: args{
				request: "/reveal?date=notAValidUnixTimestamp",
			},
			want:     revealReply{Error: ErrInvalidRequest},
			wantCode: 200,
		},
		{
			name: "error-invalid-date",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{Value: "theword"}, nil
			},
			args: args{
				request: "/reveal?date=0",
			},
			want:     revealReply{Error: ErrInvalidStartTime},
			wantCode: 200,
		},
		{
			name: "error-too-early",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{Value: "theword"}, nil
			},
			args: args{
				request: fmt.Sprintf("/reveal?date=%d", time.Now().AddDate(0, 0, 1).Unix()),
			},
			want:     revealReply{Error: ErrRevealToday},
			wantCode: 200,
		},
		{
			name: "error-getword",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{}, errors.New("ohnoes")
			},
			args: args{
				request: "/reveal",
			},
			want:     revealReply{Error: "ohnoes"},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()
			wordStore = &mockWordStore{
				mockGetForDay: tt.mockGetForDay,
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.args.request, nil)
			router.ServeHTTP(w, req)

			var response revealReply
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Error("Unable to unmarshal response body: ", w.Body.String())
			} else if w.Code != tt.wantCode {
				t.Errorf("Response Code = %d; wantCode %d", w.Code, tt.wantCode)
			} else if !reflect.DeepEqual(response, tt.want) {
				t.Errorf("Response = %#v; want %#v", response, tt.want)
			}
		})
	}
}
