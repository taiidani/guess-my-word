package data

import "github.com/google/uuid"

// User represents a single user account
type User struct {
	ID   string // Unique identifier for the user
	Name string // The user's given username
}

func newUser(username string) User {
	ret := User{Name: username}
	ret.ID = uuid.New().String()
	return ret
}
