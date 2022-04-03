package datastore

import (
	"context"
	"errors"
	"fmt"
	"guess_my_word/internal/model"
	"sort"
	"strings"
)

// MemoryClient represents an internal memory based Datastore client
type MemoryClient struct {
	Data  map[string]model.Word
	Lists map[string]model.List
}

// NewMemory instantiates a new client
func NewMemory() *MemoryClient {
	return &MemoryClient{
		Data:  map[string]model.Word{},
		Lists: map[string]model.List{},
	}
}

var errNoData = errors.New("running in local mode. Data is not available")

// GetWord will retrieve a word for the given key
func (c *MemoryClient) GetWord(ctx context.Context, key string) (model.Word, error) {
	if c.Data == nil {
		return model.Word{}, errNoData
	}

	if val, ok := c.Data[key]; ok {
		return val, nil
	}
	return model.Word{}, fmt.Errorf("word key %q not found", key)
}

// SetWord will store a Word for the given key
func (c *MemoryClient) SetWord(ctx context.Context, key string, word model.Word) error {
	if c.Data == nil {
		return errNoData
	}

	c.Data[key] = word
	return nil
}

func (c *MemoryClient) GetLists(ctx context.Context) ([]string, error) {
	if c.Lists == nil {
		return nil, errNoData
	}

	ret := []string{}
	for name := range c.Lists {
		ret = append(ret, strings.TrimPrefix(name, listCollectionPrefix))
	}

	sort.Strings(ret)
	return ret, nil
}

func (c *MemoryClient) GetList(ctx context.Context, name string) (model.List, error) {
	if c.Lists == nil {
		return model.List{}, errNoData
	}

	if val, ok := c.Lists[listCollectionPrefix+name]; ok {
		return val, nil
	}
	return model.List{}, fmt.Errorf("list key %q not found", name)
}

func (c *MemoryClient) CreateList(ctx context.Context, name string, list model.List) error {
	if c.Lists == nil {
		return errNoData
	}

	if _, err := c.GetList(ctx, name); err == nil {
		return fmt.Errorf("list already exists")
	} else if len(list.Words) == 0 {
		return fmt.Errorf("list must have at least one word")
	} else if len(list.Name) < 3 {
		return fmt.Errorf("list name is too short")
	}

	c.Lists[listCollectionPrefix+name] = list
	return nil
}

func (c *MemoryClient) DeleteList(ctx context.Context, name string) error {
	if c.Lists == nil {
		return errNoData
	}

	delete(c.Lists, listCollectionPrefix+name)
	return nil
}

func (c *MemoryClient) UpdateList(ctx context.Context, name string, list model.List) error {
	if c.Lists == nil {
		return errNoData
	}

	if _, err := c.GetList(ctx, name); err != nil {
		return fmt.Errorf("list does not exist")
	} else if len(list.Words) == 0 {
		return fmt.Errorf("list must have at least one word")
	} else if len(list.Name) < 3 {
		return fmt.Errorf("list name is too short")
	}

	c.Lists[listCollectionPrefix+name] = list
	return nil
}
