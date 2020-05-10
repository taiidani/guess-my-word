package words

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockStore struct{}

func (m *mockStore) GetWord(ctx context.Context, key string, word interface{}) error {
	return errors.New("Not happening today")
}
func (m *mockStore) SetWord(ctx context.Context, key string, word interface{}) error {
	return nil
}

func init() {
	storeClient = &mockStore{}
}

func TestGetForDay(t *testing.T) {
	type args struct {
		ctx  context.Context
		tm   time.Time
		mode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Date yesterday",
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 26, 2, 0, 0, 0, time.UTC),
				mode: "default",
			},
			want: "power",
		},
		{
			name: "Date tweak yesterday",
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 27, 2, 0, 0, 0, time.UTC).UTC().AddDate(0, 0, -1),
				mode: "default",
			},
			want: "power",
		},
		{
			name: "Unix yesterday",
			args: args{
				ctx:  context.Background(),
				tm:   time.Unix(1580083199, 0), // Sun Jan 26 23:59:59 2020 UTC
				mode: "default",
			},
			want: "power",
		},
		{
			name: "Date today",
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 27, 2, 0, 0, 0, time.UTC),
				mode: "default",
			},
			want: "tell",
		},
		{
			name: "Date tweak today",
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 26, 2, 0, 0, 0, time.UTC).UTC().AddDate(0, 0, 1),
				mode: "default",
			},
			want: "tell",
		},
		{
			name: "Unix today",
			args: args{
				ctx:  context.Background(),
				tm:   time.Unix(1580083201, 0), // Mon Jan 27 00:00:01 2020 UTC
				mode: "default",
			},
			want: "tell",
		},
		{
			name: "Hard mode date today",
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 27, 2, 0, 0, 0, time.UTC),
				mode: "hard",
			},
			want: "damans",
		},
		{
			name: "Unix OMG ERROR",
			args: args{
				ctx:  context.Background(),
				tm:   time.Unix(0, 0),
				mode: "default",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetForDay(tt.args.ctx, tt.args.tm, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetForDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetForDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	type args struct {
		ctx  context.Context
		word string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid",
			args: args{
				ctx:  context.Background(),
				word: "happy",
			},
			want: true,
		},
		{
			name: "Invalid",
			args: args{
				ctx:  context.Background(),
				word: "yppah",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.args.ctx, tt.args.word); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
