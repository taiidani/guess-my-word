package model

// Word defines information about a given key's word
type Word struct {
	Day     string  `json:"day"`     // The day that this word represents
	Guesses []Guess `json:"guesses"` // The individual guesses that were made
	Value   string  `json:"value"`
}

func (w *Word) Stats() WordStats {
	// If no one guessed that day
	if len(w.Guesses) == 0 {
		return WordStats{Word: w.Value}
	}

	ret := WordStats{
		Word:        w.Value,
		Completions: len(w.Guesses),
		BestRun:     999,
	}

	var guessCount = 0
	for _, item := range w.Guesses {
		guessCount += item.Count

		if item.Count < ret.BestRun {
			ret.BestRun = item.Count
		}
	}

	ret.AvgRun = guessCount / len(w.Guesses)
	return ret
}

type WordStats struct {
	Word        string
	Completions int
	BestRun     int
	AvgRun      int
}

type Guess struct {
	Count int `json:"count"` // The number of guesses it took to find the word
}
