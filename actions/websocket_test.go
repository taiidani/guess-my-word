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

func Test_wsHandler(t *testing.T) {
	type args struct {
		request string
	}
	tests := []struct {
		name          string
		args          args
		mockGetForDay func(ctx context.Context, tm time.Time, mode string) (model.Word, error)
		want          statsReply
		wantCode      int
	}{
		{
			name: "error-invalid-request",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{Value: "theword"}, nil
			},
			args: args{
				request: "/api/ws?date=notAValidUnixTimestamp",
			},
			want:     statsReply{Error: ErrInvalidRequest},
			wantCode: 400,
		},
		{
			name: "error-invalid-date",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{Value: "theword"}, nil
			},
			args: args{
				request: "/api/ws?date=0",
			},
			want:     statsReply{Error: ErrInvalidStartTime},
			wantCode: 400,
		},
		{
			name: "error-too-early",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{Value: "theword"}, nil
			},
			args: args{
				request: fmt.Sprintf("/api/ws?date=%d", time.Now().AddDate(0, 0, 1).Unix()),
			},
			want:     statsReply{Error: ErrRevealToday},
			wantCode: 400,
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

			var response statsReply
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

func Test_wsHandlerLoop(t *testing.T) {
	type args struct {
		ctx     context.Context
		request stats
	}
	tests := []struct {
		name          string
		args          args
		mockGetForDay func(ctx context.Context, tm time.Time, mode string) (model.Word, error)
		want          statsReply
	}{
		{
			name: "success",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{Value: "theword"}, nil
			},
			args: args{request: stats{}},
			want: statsReply{Word: model.Word{Value: "theword"}},
		},
		{
			name: "error-getword",
			mockGetForDay: func(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
				return model.Word{}, errors.New("ohnoes")
			},
			args: args{
				request: stats{},
			},
			want: statsReply{Error: "ohnoes"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wordStore = &mockWordStore{
				mockGetForDay: tt.mockGetForDay,
			}

			if tt.args.ctx == nil {
				tt.args.ctx = context.Background()
			}

			got := wsStatsHandlerLoop(tt.args.ctx, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wsStatsHandlerLoop() = %#v; want %#v", got, tt.want)
			}
		})
	}
}
