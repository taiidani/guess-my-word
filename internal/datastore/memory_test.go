package datastore

import (
	"context"
	"reflect"
	"testing"
)

func TestMemoryClient_GetWord(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		c       *MemoryClient
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			c: &MemoryClient{Data: map[string]string{
				"test": "value",
			}},
			want: "value",
		},
		{
			name: "not-found",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			c:       NewMemory(),
			wantErr: true,
		},
		{
			name: "Local mode test",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			c:       &MemoryClient{}, // No data, therefore localOnly
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c
			got, err := c.GetWord(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetWord() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("Client.GetWord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryClient_SetWord(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		word string
	}
	tests := []struct {
		name    string
		c       *MemoryClient
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				ctx:  context.Background(),
				key:  "test",
				word: "value",
			},
			c:    NewMemory(),
			want: map[string]string{"test": "value"},
		},
		{
			name: "Local mode test",
			args: args{
				ctx: context.Background(),
			},
			c:       &MemoryClient{}, // No data, therefore localOnly
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c
			if err := c.SetWord(tt.args.ctx, tt.args.key, tt.args.word); (err != nil) != tt.wantErr {
				t.Errorf("Client.SetWord() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(c.Data, tt.want) {
				t.Errorf("Client.SetWord() got = %v, want %v", c.Data, tt.want)
			}
		})
	}
}
