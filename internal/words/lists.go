package words

import (
	"context"
	"fmt"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"log"
	"strings"
)

type (
	// ListStore will manage lists
	ListStore struct {
		client Lister
	}

	// Lister represents the internal datastore for lists
	Lister interface {
		datastore.Client
	}
)

// NewListStore will return an instance of the list manager
func NewListStore(store Lister) *ListStore {
	return &ListStore{
		client: store,
	}
}

func (l *ListStore) GetLists(ctx context.Context) ([]string, error) {
	return l.client.GetLists(ctx)
}

const hardListName = "Hard"
const defaultListName = "Default"

func (l *ListStore) GetList(ctx context.Context, name string) (model.List, error) {
	switch name {
	case hardListName, strings.ToLower(hardListName):
		list, err := l.client.GetList(ctx, strings.ToLower(hardListName))
		if err != nil || len(list.Words) == 0 {
			log.Printf("hard list not present. Generating from builtin words: %s", err)
			list = model.List{
				Name:        hardListName,
				Description: "The hardest mode available -- represents the entire Scrabble dictionary",
				Color:       "422422",
				Words:       strings.Split(strings.TrimSpace(scrabbleList), "\n"),
			}
			if err := l.CreateList(ctx, hardListName, list); err != nil {
				return model.List{}, fmt.Errorf("unable to populate hard list: %w", err)
			}
		}
	case defaultListName, strings.ToLower(defaultListName):
		list, err := l.client.GetList(ctx, strings.ToLower(defaultListName))
		if err != nil || len(list.Words) == 0 {
			log.Printf("default list not present. Generating from builtin words: %s", err)
			list = model.List{
				Name:        defaultListName,
				Description: "The standard word list -- repesents about 1,000 common English words found in TV and movies",
				Words:       strings.Split(strings.TrimSpace(wordList), "\n"),
			}
			if err := l.CreateList(ctx, defaultListName, list); err != nil {
				return model.List{}, fmt.Errorf("unable to populate default list: %w", err)
			}
		}
	}

	return l.client.GetList(ctx, strings.ToLower(name))
}

func (l *ListStore) CreateList(ctx context.Context, name string, list model.List) error {
	if errs := list.Validate(); len(errs) > 0 {
		var ret error
		for _, err := range errs {
			ret = fmt.Errorf("%s: %s", ret, err)
		}
		return fmt.Errorf("validation errors: %w", ret)
	}

	return l.client.CreateList(ctx, strings.ToLower(name), list)
}

func (l *ListStore) DeleteList(ctx context.Context, name string) error {
	return l.client.DeleteList(ctx, strings.ToLower(name))
}

func (l *ListStore) UpdateList(ctx context.Context, name string, list model.List) error {
	if errs := list.Validate(); len(errs) > 0 {
		var ret error
		for _, err := range errs {
			ret = fmt.Errorf("%s: %s", ret, err)
		}
		return fmt.Errorf("validation errors: %w", ret)
	}

	return l.client.UpdateList(ctx, strings.ToLower(name), list)
}
