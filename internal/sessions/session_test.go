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
