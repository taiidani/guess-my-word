package model

// Word defines information about a given key's word
type Word struct {
	Day     string  `json:"day"`     // The day that this word represents
	Guesses []Guess `json:"guesses"` // The individual guesses that were made
	Value   string  `json:"value"`
}

type Guess struct {
	Count int `json:"count"` // The number of guesses it took to find the word
}
