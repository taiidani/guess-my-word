package words

import (
	"context"
	_ "embed"
	"fmt"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"log/slog"
	"strings"
	"time"
)

type (
	// WordStore will generate and validate words
	WordStore struct {
		client   Worder
		scrabble []string
		words    []string
	}

	// Worder represents the internal datastore for the words package
	Worder interface {
		datastore.Client
	}
)

var (
	//go:embed sowpods.txt
	scrabbleList string

	// words were taken from the original inspiration for this app, https://hryanjones.com/guess-my-word/
	// That project took the words from 1-1,000 common English words on TV and movies https://en.wiktionary.org/wiki/Wiktionary:Frequency_lists/TV/2006/1-1000
	//go:embed words.txt
	wordList string
)

// NewWordStore will return an instance of the word generator
func NewWordStore(store Worder) *WordStore {
	return &WordStore{
		client:   store,
		scrabble: strings.Split(strings.TrimSpace(scrabbleList), "\n"),
		words:    strings.Split(strings.TrimSpace(wordList), "\n"),
	}
}

// GetForDay will return a word for the given day
// This func is timezone agnostic. It will only consider the current local date
func (w *WordStore) GetForDay(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
	mode = strings.ToLower(mode)
	key := datastore.WordKey(mode, tm)
	log := slog.With("key", key)

	// Grab the word from the datastore
	word, err := w.GetWord(ctx, key)
	if err != nil {
		// Generate a new word
		log.Warn("Encountered error. Generating new word", "error", err)
		listStore := NewListStore(w.client)
		l, err := listStore.GetList(ctx, mode)
		if err != nil {
			return word, fmt.Errorf("could not get list for %q mode: %w", mode, err)
		} else if len(l.Words) == 0 {
			return word, fmt.Errorf("chosen list %q has no words", mode)
		}

		word.Value, err = w.generateWord(tm, l.Words)
		if err != nil {
			return word, err
		}

		// And store it if we're able
		log.Info("Storing generated word", "word", word)
		err = w.SetWord(ctx, key, word)
		if err != nil {
			log.Info("Encountered error storing new word", "word", word, "error", err)
		}
	}

	if word.Day == "" {
		word.Day = tm.Format("2006-01-02")
	}

	return word, nil
}

// Validate will confirm if a given word is valid
func (w *WordStore) Validate(ctx context.Context, word string) bool {
	for _, line := range w.scrabble {
		if line == word {
			return true
		}
	}

	return false
}

func (w *WordStore) GetWord(ctx context.Context, key string) (model.Word, error) {
	return w.client.GetWord(ctx, strings.ToLower(key))
}

func (w *WordStore) SetWord(ctx context.Context, key string, word model.Word) error {
	return w.client.SetWord(ctx, strings.ToLower(key), word)
}

func (w *WordStore) generateWord(seed time.Time, words []string) (string, error) {
	if seed.Unix() == 0 {
		return "", fmt.Errorf("invalid timestamp for word")
	}
	return words[(seed.Year()*seed.YearDay())%len(words)], nil
}
