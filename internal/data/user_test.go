package data

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func Test_newUser(t *testing.T) {
	// Force the UUID seed to a known value
	seed := bytes.NewBufferString("12345678901233456789")
	uuid.SetRand(seed)
	defer uuid.SetRand(nil)

	// Aaaand test
	tests := []struct {
		name     string
		username string
		want     User
	}{
		{
			name:     "New user",
			username: "Flynn",
			want:     User{ID: "31323334-3536-4738-b930-313233333435", Name: "Flynn"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUser(tt.username); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
