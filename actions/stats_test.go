package actions

import (
	"context"
	"errors"
	"guess_my_word/internal/model"
	"reflect"
	"testing"
	"time"
)

func Test_refreshStats(t *testing.T) {
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

			got := refreshStats(tt.args.ctx, tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("refreshStats() = %#v; want %#v", got, tt.want)
			}
		})
	}
}
