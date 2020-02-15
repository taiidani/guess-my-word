package data

import (
	"reflect"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v2"
)

func TestLoadDate(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantDte *Date
		wantErr bool
	}{
		{
			name: "Valid Date",
			args: args{
				time.Date(2010, time.January, 10, 23, 00, 00, 00, time.UTC),
			},
			wantDte: &Date{
				ID: "2010-01-10",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := badger.DefaultOptions("")
			opts.InMemory = true
			db, teardown := NewBadgerBackend(&opts)
			SetBackend(db)
			defer teardown()

			// Save the date first
			if err := tt.wantDte.Save(); err != nil {
				t.Errorf("Could not save date: %v", err)
			}

			// Then load the date
			gotDte, err := LoadDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDte, tt.wantDte) {
				t.Errorf("LoadDate() = %v, want %v", gotDte, tt.wantDte)
			}
		})
	}
}
