package words

import (
	"context"
	_ "embed"
	"fmt"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"log"
	"strings"
	"time"
)

type (
	// WordStore will generate and validate words
	WordStore struct {
		storeClient Store
		scrabble    []string
		words       []string
	}

	// Store represents the internal datastore for the words package
	Store interface {
		GetWord(ctx context.Context, key string) (model.Word, error)
		SetWord(ctx context.Context, key string, word model.Word) error
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
func NewWordStore(store Store) *WordStore {
	return &WordStore{
		storeClient: store,
		scrabble:    strings.Split(strings.TrimSpace(scrabbleList), "\n"),
		words:       strings.Split(strings.TrimSpace(wordList), "\n"),
	}
}

// GetForDay will return a word for the given day
// This func is timezone agnostic. It will only consider the current local date
func (w *WordStore) GetForDay(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
	key := datastore.WordKey(mode, tm)
	log.Println("Getting word for day at", key)

	// Grab the word from the datastore
	word, err := w.storeClient.GetWord(ctx, key)
	if err != nil {
		// Generate a new word
		log.Printf("Encountered error '%s'. Generating new word for key '%s'", err, key)
		word.Value, err = w.generateWord(tm, w.getWordList(mode))
		if err != nil {
			return word, err
		}

		// And store it if we're able
		log.Printf("Storing generated word '%v' at key '%s'", word, key)
		err = w.storeClient.SetWord(ctx, key, word)
		if err != nil {
			log.Printf("Encountered error storing new word '%v' at key '%s': %s", word, key, err)
		}
	}

	if word.Day == "" {
		word.Day = tm.Format("2006-01-02")
	}

	return word, err
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

func (w *WordStore) generateWord(seed time.Time, words []string) (string, error) {
	if seed.Unix() == 0 {
		return "", fmt.Errorf("invalid timestamp for word")
	}
	return words[(seed.Year()*seed.YearDay())%len(words)], nil
}

func (w *WordStore) getWordList(mode string) []string {
	switch mode {
	case "hard":
		return w.scrabble
	default:
		return w.words
	}
}
