package data

import "time"

// Leaderboard tracks the Guess winners for a given day
type Leaderboard struct {
	Entries []LeaderboardEntry // The list of entries on the leaderboard
}

// LeaderboardEntry tracks a single user on a leaderboard
type LeaderboardEntry struct {
	User    User          // The name of the user
	Guesses int64         // The number of guesses it took for the user to guess the word
	Time    time.Duration // The length of time it took for the user to guess the word
}
