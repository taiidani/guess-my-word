package data

import (
	"reflect"
	"testing"

	"github.com/dgraph-io/badger/v2"
)

func TestBadgerBackend_SetGet(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				key:   "lorem",
				value: "ipsum",
			},
			want: []byte("ipsum"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := badger.DefaultOptions("")
			opts.InMemory = true
			b, teardown := NewBadgerBackend(&opts)
			defer teardown()

			if err := b.Set(tt.args.key, []byte(tt.args.value)); err != nil {
				t.Errorf("Failed during set: %v", err)
			}

			got, err := b.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("BadgerBackend.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BadgerBackend.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
