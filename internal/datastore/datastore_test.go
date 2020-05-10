package datastore

import (
	"context"
	"testing"
)

func init() {
	LocalOnly = true
}

func TestClient_GetWord(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		word interface{}
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "Local mode test",
			args: args{
				ctx: context.Background(),
				key: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			if err := c.GetWord(tt.args.ctx, tt.args.key, tt.args.word); (err != nil) != tt.wantErr {
				t.Errorf("Client.GetWord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SetWord(t *testing.T) {
	type args struct {
		ctx  context.Context
		key  string
		word interface{}
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "Local mode test",
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			if err := c.SetWord(tt.args.ctx, tt.args.key, tt.args.word); (err != nil) != tt.wantErr {
				t.Errorf("Client.SetWord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
