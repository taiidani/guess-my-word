package actions

import (
	"testing"
	"time"
)

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
