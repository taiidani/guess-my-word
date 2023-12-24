package words

import (
	"testing"
)

func TestDictionary_Validate(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name      string
		args      args
		wantI     int
		wantFound bool
	}{
		{
			name: "Valid",
			args: args{
				word: "happy",
			},
			wantI:     101650,
			wantFound: true,
		},
		{
			name: "Invalid",
			args: args{
				word: "yppah",
			},
			wantI:     0,
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Dictionary{
				Words: ScrabbleDictionary.Words,
			}

			gotI, gotFound := w.Validate(tt.args.word)
			if gotI != tt.wantI {
				t.Errorf("Validate() i = %d, want %d", gotI, tt.wantI)
			}
			if gotFound != tt.wantFound {
				t.Errorf("Validate() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
