package data

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"encoding/json"
)

type (
	// Date stores all data for a given day
	Date struct {
		ID          string       // The identifier for this date
		Word        string       // The word of the day
		Leaderboard Leaderboard  // The leaderboard for the day
		Suggestions []Suggestion // The list of suggested words for the following day
	}
)

const dateFormat = "2006-01-02"

// NewDate will generate a new Date object for the given day
// It is general best practice to first attempt to LoadDate before generating a new one
func NewDate(date time.Time) *Date {
	id := date.Format(dateFormat)
	word, err := generateWord(date)
	if err != nil {
		log.Panicf("Unable to generate new date for %s: %s", id, err)
	}

	return &Date{
		ID:          id,
		Word:        word,
		Leaderboard: Leaderboard{},
		Suggestions: []Suggestion{},
	}
}

// LoadDate will attempt to load the given date from the data store
func LoadDate(date time.Time) (dte *Date, err error) {
	data, err := db.Get(date.Format(dateFormat))
	if err != nil {
		return nil, fmt.Errorf("Could not load date for %s: %w", date.Format(dateFormat), err)
	}

	dte = &Date{}
	err = json.Unmarshal(data, dte)
	return
}

// Save will save the current Date to the backend
func (d *Date) Save() error {
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err := enc.Encode(d); err != nil {
		return err
	}

	err := db.Set(d.ID, data.Bytes())
	if err != nil {
		return fmt.Errorf("Unable to save date with ID %s: %w", d.ID, err)
	}

	return nil
}
