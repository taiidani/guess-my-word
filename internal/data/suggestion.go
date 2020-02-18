package data

import "strings"

// Suggestion tracks a word being suggested by users for a subsequent day
type Suggestion struct {
	Word  string // The word being suggested
	Users []User // The list of users that have suggested it
}

// AddSuggestion will record a user's suggested word for tomorrow
func (d *Date) AddSuggestion(word string, user User) {
	var suggestion *Suggestion
	for i, existingSuggestion := range d.Suggestions {
		if strings.ToLower(existingSuggestion.Word) == strings.ToLower(word) {
			suggestion = &d.Suggestions[i]
			break
		}
	}

	if suggestion == nil {
		// No one has suggested this before!
		suggestion = &Suggestion{
			Word:  word,
			Users: []User{user},
		}
		d.Suggestions = append(d.Suggestions, *suggestion)
		return
	}

	for _, existingUser := range suggestion.Users {
		if existingUser.ID == user.ID {
			// User already exists. Can't suggest twice!
			return
		}
	}

	// A brand new user making the suggestion!
	suggestion.Users = append(suggestion.Users, user)
}
