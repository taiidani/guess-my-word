package datastore

import "time"

const (
	wordCollectionPrefix = "words/"
)

// WordKey returns the primary key used to store the given mode/day word data
func WordKey(mode string, tm time.Time) string {
	return mode + "/day/" + tm.Format("2006-01-02")
}
