package data

import (
	"fmt"
	"time"
)

func generateWord(seed time.Time) (string, error) {
	if seed.Unix() == 0 {
		return "", fmt.Errorf("Invalid timestamp for word")
	}

	day := seed.UTC()
	return words[(day.Year()*day.YearDay())%len(words)], nil
}

// ValidateWord validates that the given word is a valid Scrabble word
func ValidateWord(word string) bool {
	for _, line := range scrabble {
		if line == word {
			return true
		}
	}

	return false
}
