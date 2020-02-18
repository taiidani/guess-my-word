package data

import (
	"reflect"
	"testing"
)

func TestDate_AddSuggestion(t *testing.T) {
	type args struct {
		word string
		user User
	}
	tests := []struct {
		name            string
		suggestions     []Suggestion
		args            args
		wantSuggestions []Suggestion
	}{
		{
			name:        "New Suggestion, New User",
			suggestions: []Suggestion{{Word: "lobster", Users: []User{}}},
			args: args{
				word: "gopher",
				user: User{ID: "1234abcd", Name: "Flynn"},
			},
			wantSuggestions: []Suggestion{
				{Word: "lobster", Users: []User{}},
				{Word: "gopher", Users: []User{{ID: "1234abcd", Name: "Flynn"}}},
			},
		},
		{
			name:        "Existing Suggestion, New User",
			suggestions: []Suggestion{{Word: "gopher", Users: []User{}}},
			args: args{
				word: "gopher",
				user: User{ID: "1234abcd", Name: "Flynn"},
			},
			wantSuggestions: []Suggestion{
				{Word: "gopher", Users: []User{{ID: "1234abcd", Name: "Flynn"}}},
			},
		},
		{
			name:        "Existing Suggestion, Existing User",
			suggestions: []Suggestion{{Word: "gopher", Users: []User{{ID: "1234abcd", Name: "Flynn"}}}},
			args: args{
				word: "gopher",
				user: User{ID: "1234abcd", Name: "Flynn"},
			},
			wantSuggestions: []Suggestion{
				{Word: "gopher", Users: []User{{ID: "1234abcd", Name: "Flynn"}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Date{Suggestions: tt.suggestions}
			d.AddSuggestion(tt.args.word, tt.args.user)

			if !reflect.DeepEqual(d.Suggestions, tt.wantSuggestions) {
				t.Errorf("[]suggestions = %v, want %v", d.Suggestions, tt.wantSuggestions)
			}
		})
	}
}
