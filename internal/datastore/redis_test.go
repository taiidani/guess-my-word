package datastore

import (
	"context"
	"guess_my_word/internal/model"
	"reflect"
	"testing"
)

func TestRedisClient_GetWord(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		c       *RedisClient
		args    args
		want    model.Word
		wantErr bool
	}{
		{
			name: "Local mode test",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			c:       &RedisClient{}, // No client, therefore localOnly
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetWord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisClient_SetWord(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		word model.Word
	}
	tests := []struct {
		name    string
		c       *RedisClient
		args    args
		wantErr bool
	}{
		{
			name: "Local mode test",
			args: args{
				ctx: context.Background(),
			},
			c:       &RedisClient{}, // No client, therefore localOnly
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c
			if err := c.SetWord(tt.args.ctx, tt.args.key, tt.args.word); (err != nil) != tt.wantErr {
				t.Errorf("Client.SetWord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
