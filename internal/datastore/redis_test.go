package datastore

import (
	"context"
	"fmt"
	"guess_my_word/internal/model"
	"reflect"
	"testing"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type mockRedisClient struct {
	mockGet    func(ctx context.Context, key string) *redis.StringCmd
	mockSet    func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	mockKeys   func(ctx context.Context, pattern string) *redis.StringSliceCmd
	mockExists func(ctx context.Context, keys ...string) *redis.IntCmd
	mockDel    func(ctx context.Context, keys ...string) *redis.IntCmd
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return m.mockGet(ctx, key)
}
func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.mockSet(ctx, key, value, expiration)
}
func (m *mockRedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return m.mockKeys(ctx, pattern)
}
func (m *mockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return m.mockExists(ctx, keys...)
}
func (m *mockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return m.mockDel(ctx, keys...)
}

func TestRedisClient_GetWord(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		c       redisClient
		args    args
		want    model.Word
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			c: &mockRedisClient{
				mockGet: func(ctx context.Context, key string) *redis.StringCmd {
					return redis.NewStringResult(`{ "day": "2021-01-01" }`, nil)
				},
			},
			want: model.Word{Day: "2021-01-01"},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			c: &mockRedisClient{
				mockGet: func(ctx context.Context, key string) *redis.StringCmd {
					return redis.NewStringResult(``, fmt.Errorf("oh noez"))
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := RedisClient{
				client: tt.c,
			}
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
		c       redisClient
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				ctx: context.Background(),
			},
			c: &mockRedisClient{
				mockSet: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
					return redis.NewStatusResult("OK", nil)
				},
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
			},
			c: &mockRedisClient{
				mockSet: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
					return redis.NewStatusResult("", fmt.Errorf("oh noez"))
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := RedisClient{
				client: tt.c,
			}
			if err := c.SetWord(tt.args.ctx, tt.args.key, tt.args.word); (err != nil) != tt.wantErr {
				t.Errorf("Client.SetWord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
