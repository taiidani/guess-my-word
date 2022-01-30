package datastore

import (
	"context"
	"errors"
	"fmt"
	"guess_my_word/internal/model"
)

// MemoryClient represents an internal memory based Datastore client
type MemoryClient struct {
	Data map[string]model.Word
}

// NewMemory instantiates a new client
func NewMemory() *MemoryClient {
	return &MemoryClient{Data: map[string]model.Word{}}
}

// GetWord will retrieve a word for the given key
func (c *MemoryClient) GetWord(ctx context.Context, key string) (model.Word, error) {
	if c.Data == nil {
		return model.Word{}, errors.New("running in local mode. Data has not been retrieved")
	}

	if val, ok := c.Data[key]; ok {
		return val, nil
	}
	return model.Word{}, fmt.Errorf("key not found")
}

// SetWord will store a Word for the given key
func (c *MemoryClient) SetWord(ctx context.Context, key string, word model.Word) error {
	if c.Data == nil {
		return fmt.Errorf("running in local mode. Data has not been stored")
	}

	c.Data[key] = word
	return nil
}
