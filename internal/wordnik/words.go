package wordnik

import (
	"context"
	"net/url"
	"strconv"
)

type (
	// Words represents the Wordnik Words API
	// https://developer.wordnik.com/docs#/words
	Words interface {
		RandomWord(context.Context, RandomWordRequest) (RandomWordResponse, error)
	}

	words struct {
		*Client
	}

	// RandomWordRequest contains the request parameters for the RandomWord endpoint
	RandomWordRequest struct {
		MinLength      uint64
		MinCorpusCount uint64
	}

	// RandomWordResponse contains the response for the RandomWord endpoint
	RandomWordResponse struct {
		Word string
	}
)

func (c *Client) newWords() *words {
	return &words{Client: c}
}

func (w *words) RandomWord(ctx context.Context, req RandomWordRequest) (resp RandomWordResponse, err error) {
	params := url.Values{}
	params.Add("minLength", strconv.FormatUint(req.MinLength, 10))
	params.Add("minCorpusCount", strconv.FormatUint(req.MinCorpusCount, 10))

	r, err := w.get(ctx, "words.json/randomWord", params)
	if err != nil {
		return resp, err
	}

	err = w.handleRequest(r, &resp)
	return resp, err
}
