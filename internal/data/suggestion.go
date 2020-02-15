package data

// Suggestion tracks a word being suggested by users for a subsequent day
type Suggestion struct {
	Word  string // The word being suggested
	Users []User // The list of users that have suggested it
}
