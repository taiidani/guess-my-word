package words

import (
	"context"
	"guess_my_word/internal/wordnik"
	"os"
)

var (
	// wordClient contains the internal Wordnik API client
	wordClient *wordnik.Client
)

func init() {
	wordClient = wordnik.NewClient(os.Getenv("WORDNIK_API_KEY"))
}

func generateWordnikWord(ctx context.Context) (string, error) {
	resp, err := wordClient.Words.RandomWord(ctx, wordnik.RandomWordRequest{
		MinLength:      4,
		MinCorpusCount: 700000,
	})
	return resp.Word, err
}
