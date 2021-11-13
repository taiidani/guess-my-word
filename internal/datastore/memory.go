package datastore

import (
	"context"
	"errors"
	"fmt"
)

// MemoryClient represents an internal memory based Datastore client
type MemoryClient struct {
	Data map[string]string
}

// NewMemory instantiates a new client
func NewMemory() *MemoryClient {
	return &MemoryClient{Data: map[string]string{}}
}

// GetWord will retrieve a word for the given key
func (c *MemoryClient) GetWord(ctx context.Context, key string) (string, error) {
	if c.Data == nil {
		return "", errors.New("running in local mode. Data has not been retrieved")
	}

	if val, ok := c.Data[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

// SetWord will store a Word for the given key
func (c *MemoryClient) SetWord(ctx context.Context, key string, word string) error {
	if c.Data == nil {
		return fmt.Errorf("Running in local mode. Data has not been stored")
	}

	c.Data[key] = word
	return nil
}
