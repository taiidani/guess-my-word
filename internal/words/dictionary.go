package words

import (
	_ "embed"
	"strings"
)

type Dictionary struct {
	Words []string
}

//go:embed sowpods.txt
var dictionaryList string

var ScrabbleDictionary = &Dictionary{
	Words: strings.Split(strings.TrimSpace(dictionaryList), "\n"),
}

// Validate will confirm if a given word is valid.
//
// It will return the position of the word in the Scrabble list as well as
// a boolean indicating if it was found.
func (d *Dictionary) Validate(word string) (int, bool) {
	for i, line := range d.Words {
		if line == word {
			return i, true
		}
	}

	return 0, false
}

// DictionarySize will return the total size of the Scrabble dictionary.
func (d *Dictionary) Size() int {
	return len(d.Words)
}
