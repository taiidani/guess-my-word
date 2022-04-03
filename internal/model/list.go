package model

import (
	"fmt"
	"regexp"
)

// List defines a list of guessable words
type List struct {
	Name        string   `json:"name"`               // The given name of the list
	Description string   `json:"description"`        // A description of the list's contents
	Color       string   `json:"color,omitempty"`    // The background color of the app when using this list
	Password    string   `json:"password,omitempty"` // SHA1 hash of the list author's password
	Words       []string `json:"words"`              // The words in the list
}

var (
	listNameRegex  = regexp.MustCompile("^[A-Za-z ]{4,}$")
	listColorRegex = regexp.MustCompile("^[a-f0-9]{6}$")
	listWordRegex  = regexp.MustCompile("^[A-Za-z]{2,}$")
)

func (l *List) Validate() []error {
	ret := []error{}
	if !listNameRegex.MatchString(l.Name) {
		ret = append(ret, fmt.Errorf("invalid name. Names must be at least 4 letters and contain only letters and spaces"))
	}

	if len(l.Color) > 0 && !listColorRegex.MatchString(l.Color) {
		ret = append(ret, fmt.Errorf("invalid color. Colors must match a CSS color code such as 123abc"))
	}

	if len(l.Words) < 3 {
		ret = append(ret, fmt.Errorf("all lists must contain at least 3 words"))
	}

	for _, word := range l.Words {
		if !listWordRegex.MatchString(word) {
			ret = append(ret, fmt.Errorf("invalid word: %q. Words be at least two letters and contain only letters", word))
		}
	}

	return ret
}
