package wordnik

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_words_RandomWord(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		ctx context.Context
		req RandomWordRequest
	}
	tests := []struct {
		name    string
		handler http.HandlerFunc
		fields  fields
		args    args
		want    RandomWordResponse
		wantErr bool
	}{
		{
			name: "Success",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				fmt.Fprintln(w, `{ "word": "Nope" }`)
			}),
			fields: fields{Client: &Client{
				apiKey: apiKeyTODO,
				client: http.DefaultClient,
			}},
			args: args{
				ctx: context.Background(),
				req: RandomWordRequest{
					MinLength:      4,
					MinCorpusCount: 700000,
				},
			},
			want: RandomWordResponse{
				Word: "Nope",
			},
		},
		{
			name: "Rate Limited",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(429)
				fmt.Fprintln(w, `{ "message": "Rate Limit Reached" }`)
			}),
			fields: fields{Client: &Client{
				apiKey: apiKeyTODO,
				client: http.DefaultClient,
			}},
			args: args{
				ctx: context.Background(),
				req: RandomWordRequest{
					MinLength:      4,
					MinCorpusCount: 700000,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				fmt.Fprintln(w, `Not valid json`)
			}),
			fields: fields{Client: &Client{
				apiKey: apiKeyTODO,
				client: http.DefaultClient,
			}},
			args: args{
				ctx: context.Background(),
				req: RandomWordRequest{
					MinLength:      4,
					MinCorpusCount: 700000,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()
			urlPrefix = ts.URL

			w := &words{
				Client: tt.fields.Client,
			}
			got, err := w.RandomWord(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("words.RandomWord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("words.RandomWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
