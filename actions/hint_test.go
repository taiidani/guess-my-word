package actions

import "testing"

func Test_getWordHint(t *testing.T) {
	type args struct {
		h    hint
		word string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1 character",
			args: args{
				h: hint{
					After:  "zoo",
					Before: "train",
				},
				word: "wave",
			},
			want: "w",
		},
		{
			name: "2 characters",
			args: args{
				h: hint{
					After:  "tear",
					Before: "train",
				},
				word: "time",
			},
			want: "ti",
		},
		{
			name: "Almost there",
			args: args{
				h: hint{
					After:  "tray",
					Before: "trays",
				},
				word: "traybit",
			},
			want: "trayb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWordHint(tt.args.h, tt.args.word); got != tt.want {
				t.Errorf("getWordHint() = %v, want %v", got, tt.want)
			}
		})
	}
}
