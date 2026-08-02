package hangman

import (
	"context"
	"slices"
	"strings"
)

type Session struct {
	word    string
	letters []rune
}

const missingLetter rune = '*'

func (s *Session) updateGuess(ctx context.Context) string {
	newGuess := strings.Builder{}

	for _, r := range s.word {
		if slices.Contains(s.letters, r) {
			newGuess.WriteRune(r)
		} else {
			newGuess.WriteRune(missingLetter)
		}
	}

	return newGuess.String()
}
