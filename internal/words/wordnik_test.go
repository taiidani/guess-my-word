package words

import (
	"context"
	"errors"
	"guess_my_word/internal/wordnik"
	"testing"
)

type mockWordsService struct {
	mockRandomWord func(ctx context.Context, req wordnik.RandomWordRequest) (resp wordnik.RandomWordResponse, err error)
}

func (m *mockWordsService) RandomWord(ctx context.Context, req wordnik.RandomWordRequest) (resp wordnik.RandomWordResponse, err error) {
	return m.mockRandomWord(ctx, req)
}

func Test_generateWordnikWord(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		words   *mockWordsService
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			words: &mockWordsService{
				mockRandomWord: func(ctx context.Context, req wordnik.RandomWordRequest) (resp wordnik.RandomWordResponse, err error) {
					return wordnik.RandomWordResponse{Word: "hi"}, nil
				},
			},
			args: args{ctx: context.Background()},
			want: "hi",
		},
		{
			name: "Failure",
			words: &mockWordsService{
				mockRandomWord: func(ctx context.Context, req wordnik.RandomWordRequest) (resp wordnik.RandomWordResponse, err error) {
					return wordnik.RandomWordResponse{Word: ""}, errors.New("oh noez")
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wordClient.Words = tt.words

			got, err := generateWordnikWord(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateWordnikWord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateWordnikWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
