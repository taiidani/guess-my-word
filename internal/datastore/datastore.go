package datastore

import (
	"context"
	"guess_my_word/internal/model"
	"time"
)

const (
	listCollectionPrefix = "lists/"
	wordCollectionPrefix = "words/"
)

type Client interface {
	GetLists(ctx context.Context) ([]string, error)
	GetList(ctx context.Context, name string) (model.List, error)
	CreateList(ctx context.Context, name string, list model.List) error
	DeleteList(ctx context.Context, name string) error
	UpdateList(ctx context.Context, name string, list model.List) error
	GetWord(ctx context.Context, key string) (model.Word, error)
	SetWord(ctx context.Context, key string, word model.Word) error
}

// WordKey returns the primary key used to store the given mode/day word data
func WordKey(mode string, tm time.Time) string {
	return mode + "/day/" + tm.Format("2006-01-02")
}
