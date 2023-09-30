package words

import (
	"context"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewWordStore(t *testing.T) {
	w := NewWordStore(nil)
	if w == nil {
		t.Error("Received nil for word instance")
	}
}

func TestWordStore_GetForDay(t *testing.T) {
	type fields struct {
		client Worder
	}
	type args struct {
		ctx  context.Context
		tm   time.Time
		mode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Word
		wantErr bool
	}{
		{
			name:   "Date yesterday",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 26, 2, 0, 0, 0, time.UTC),
				mode: defaultListName,
			},
			want: model.Word{Day: "2020-01-26", Value: "power"},
		},
		{
			name:   "Date tweak yesterday",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 27, 2, 0, 0, 0, time.UTC).UTC().AddDate(0, 0, -1),
				mode: defaultListName,
			},
			want: model.Word{Day: "2020-01-26", Value: "power"},
		},
		{
			name:   "Unix yesterday",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx:  context.Background(),
				tm:   time.Unix(1580083199, 0), // Sun Jan 26 23:59:59 2020 UTC
				mode: defaultListName,
			},
			want: model.Word{Day: "2020-01-26", Value: "power"},
		},
		{
			name:   "Date today",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 27, 2, 0, 0, 0, time.UTC),
				mode: defaultListName,
			},
			want: model.Word{Day: "2020-01-27", Value: "tell"},
		},
		{
			name:   "Date today TZ",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx: context.Background(),
				tm: time.Date(2020, time.January, 27, 2, 0, 0, 0, func() *time.Location {
					ret, err := time.LoadLocation("America/Los_Angeles")
					if err != nil {
						t.Fatal("Test machine timezone data not populated")
					}
					return ret
				}()),
				mode: defaultListName,
			},
			want: model.Word{Day: "2020-01-27", Value: "tell"},
		},
		{
			name:   "Date tweak today",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 26, 2, 0, 0, 0, time.UTC).UTC().AddDate(0, 0, 1),
				mode: defaultListName,
			},
			want: model.Word{Day: "2020-01-27", Value: "tell"},
		},
		{
			name:   "Unix today",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx:  context.Background(),
				tm:   time.Unix(1580083201, 0).UTC(), // Mon Jan 27 00:00:01 2020 UTC
				mode: defaultListName,
			},
			want: model.Word{Day: "2020-01-27", Value: "tell"},
		},
		{
			name:   "Hard mode date today",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx:  context.Background(),
				tm:   time.Date(2020, time.January, 27, 2, 0, 0, 0, time.UTC),
				mode: hardListName,
			},
			want: model.Word{Day: "2020-01-27", Value: "damans"},
		},
		{
			name:   "Unix OMG ERROR",
			fields: fields{client: datastore.NewMemory()},
			args: args{
				ctx:  context.Background(),
				tm:   time.Unix(0, 0),
				mode: defaultListName,
			},
			want:    model.Word{Value: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WordStore{
				client:   tt.fields.client,
				scrabble: strings.Split(strings.TrimSpace(scrabbleList), "\n"),
				words:    strings.Split(strings.TrimSpace(wordList), "\n"),
			}

			got, err := w.GetForDay(tt.args.ctx, tt.args.tm, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetForDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetForDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordStore_Validate(t *testing.T) {
	type args struct {
		ctx  context.Context
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
				ctx:  context.Background(),
				word: "happy",
			},
			wantI:     101650,
			wantFound: true,
		},
		{
			name: "Invalid",
			args: args{
				ctx:  context.Background(),
				word: "yppah",
			},
			wantI:     0,
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WordStore{
				scrabble: strings.Split(strings.TrimSpace(scrabbleList), "\n"),
			}

			gotI, gotFound := w.Validate(tt.args.ctx, tt.args.word)
			if gotI != tt.wantI {
				t.Errorf("Validate() i = %d, want %d", gotI, tt.wantI)
			}
			if gotFound != tt.wantFound {
				t.Errorf("Validate() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
