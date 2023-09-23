package sessions

import (
	"testing"
	"time"
)

func TestSessionMode_CommonGuessPrefix(t *testing.T) {
	type fields struct {
		Start  time.Time
		End    *time.Time
		Before []string
		After  []string
		Answer string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "tree",
			fields: fields{
				Before: []string{"apple", "trace"},
				After:  []string{"tree", "yarn"},
			},
			want: "tr",
		},
		{
			name: "cash",
			fields: fields{
				Before: []string{"apple", "cashier"},
				After:  []string{"cashiers", "tree"},
			},
			want: "cashier",
		},
		{
			name: "no common items",
			fields: fields{
				Before: []string{"apple"},
				After:  []string{"tree"},
			},
			want: "",
		},
		{
			name:   "empty",
			fields: fields{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SessionMode{
				Start:  tt.fields.Start,
				End:    tt.fields.End,
				Before: tt.fields.Before,
				After:  tt.fields.After,
				Answer: tt.fields.Answer,
			}
			if got := m.CommonGuessPrefix(); got != tt.want {
				t.Errorf("SessionMode.CommonGuessPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionMode_Stale(t *testing.T) {
	type fields struct {
		Start time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "not stale",
			fields: fields{Start: time.Now()},
			want:   false,
		},
		{
			name:   "stale",
			fields: fields{Start: time.Now().Add(time.Hour * -24)},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SessionMode{
				Start: tt.fields.Start,
			}
			if got := m.Stale(); got != tt.want {
				t.Errorf("SessionMode.Stale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionMode_RemainingTime(t *testing.T) {
	tests := []struct {
		name              string
		remainingSeedTime time.Time
		want              string
	}{
		{
			name:              "24 hours",
			remainingSeedTime: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			want:              "24 hours, 0 minutes",
		},
		{
			name:              "6 hours, 30 minutes",
			remainingSeedTime: time.Date(2020, time.January, 1, 17, 30, 0, 0, time.UTC),
			want:              "6 hours, 30 minutes",
		},
		{
			name:              "7 hours, 59 minutes",
			remainingSeedTime: time.Date(2020, time.January, 1, 16, 1, 0, 0, time.UTC),
			want:              "7 hours, 59 minutes",
		},
		{
			name:              "1 hour, 1 minute",
			remainingSeedTime: time.Date(2020, time.January, 1, 22, 59, 0, 0, time.UTC),
			want:              "1 hour, 1 minute",
		},
		{
			name:              "0 hours, 1 minute",
			remainingSeedTime: time.Date(2020, time.January, 1, 23, 59, 0, 0, time.UTC),
			want:              "0 hours, 1 minute",
		},
		{
			name:              "Next month",
			remainingSeedTime: time.Date(2020, time.January, 31, 0, 0, 0, 0, time.UTC),
			want:              "24 hours, 0 minutes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SessionMode{}
			remainingSeedTime = tt.remainingSeedTime

			if got := m.RemainingTime(); got != tt.want {
				t.Errorf("SessionMode.RemainingTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
