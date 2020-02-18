package data

import (
	"errors"
	"time"
)

// Leaderboard tracks the Guess winners for a given day
type Leaderboard struct {
	Entries []LeaderboardEntry // The list of entries on the leaderboard
}

// LeaderboardEntry tracks a single user on a leaderboard
type LeaderboardEntry struct {
	User    User          // The name of the user
	Guesses uint64        // The number of guesses it took for the user to guess the word
	Time    time.Duration // The length of time it took for the user to guess the word
}

// AddLeaderboard will add the given user information to the day's leaderboard
func (d *Date) AddLeaderboard(username string, guesses uint64, duration time.Duration) (*LeaderboardEntry, error) {
	for i, existingEntry := range d.Leaderboard.Entries {
		// The user is already on the board!
		if existingEntry.User.Name == username {
			return &d.Leaderboard.Entries[i], errors.New("Entry already exists in leaderboard")
		}
	}

	// New user!
	user := newUser(username)
	entry := LeaderboardEntry{
		User:    user,
		Guesses: guesses,
		Time:    duration,
	}
	d.Leaderboard.Entries = append(d.Leaderboard.Entries, entry)
	return &d.Leaderboard.Entries[len(d.Leaderboard.Entries)-1], nil
}
