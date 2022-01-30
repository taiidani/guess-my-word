package model

// Word defines information about a given key's word
type Word struct {
	Guesses []Guess // The individual guesses that were made
	Value   string
}

type Guess struct {
	Count int // The number of guesses it took to find the word
}
