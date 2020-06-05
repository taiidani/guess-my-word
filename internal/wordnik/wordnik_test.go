package wordnik

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("key")
	if client == nil || client.client == nil {
		t.Error("Client object was nil")
	} else if client.apiKey != "key" {
		t.Errorf("Wanted 'key' for apiKey, got '%s'", client.apiKey)
	}
}
