package datastore

import (
	"context"
	"errors"
	"fmt"
	"guess_my_word/internal/model"
	"sort"
	"strings"
	"sync"
)

// MemoryClient represents an internal memory based Datastore client
type MemoryClient struct {
	data     map[string]model.Word
	lists    map[string]model.List
	mapMutex sync.Mutex
}

// NewMemory instantiates a new client
func NewMemory() *MemoryClient {
	return &MemoryClient{
		data:     map[string]model.Word{},
		lists:    map[string]model.List{},
		mapMutex: sync.Mutex{},
	}
}

var errNoData = errors.New("running in local mode. Data is not available")

// GetWord will retrieve a word for the given key
func (c *MemoryClient) GetWord(ctx context.Context, key string) (model.Word, error) {
	if c.data == nil {
		return model.Word{}, errNoData
	}

	if val, ok := c.data[key]; ok {
		return val, nil
	}
	return model.Word{}, fmt.Errorf("word key %q not found", key)
}

// SetWord will store a Word for the given key
func (c *MemoryClient) SetWord(ctx context.Context, key string, word model.Word) error {
	c.mapMutex.Lock()
	defer c.mapMutex.Unlock()

	if c.data == nil {
		return errNoData
	}

	c.data[key] = word
	return nil
}

func (c *MemoryClient) GetLists(ctx context.Context) ([]string, error) {
	if c.lists == nil {
		return nil, errNoData
	}

	ret := []string{}
	for name := range c.lists {
		ret = append(ret, strings.TrimPrefix(name, listCollectionPrefix))
	}

	sort.Strings(ret)
	return ret, nil
}

func (c *MemoryClient) GetList(ctx context.Context, name string) (model.List, error) {
	if c.lists == nil {
		return model.List{}, errNoData
	}

	if val, ok := c.lists[listCollectionPrefix+name]; ok {
		return val, nil
	}
	return model.List{}, fmt.Errorf("list key %q not found", name)
}

func (c *MemoryClient) CreateList(ctx context.Context, name string, list model.List) error {
	c.mapMutex.Lock()
	defer c.mapMutex.Unlock()

	if c.lists == nil {
		return errNoData
	}

	if _, err := c.GetList(ctx, name); err == nil {
		return fmt.Errorf("list already exists")
	} else if len(list.Words) == 0 {
		return fmt.Errorf("list must have at least one word")
	} else if len(list.Name) < 3 {
		return fmt.Errorf("list name is too short")
	}

	c.lists[listCollectionPrefix+name] = list
	return nil
}

func (c *MemoryClient) DeleteList(ctx context.Context, name string) error {
	if c.lists == nil {
		return errNoData
	}

	delete(c.lists, listCollectionPrefix+name)
	return nil
}

func (c *MemoryClient) UpdateList(ctx context.Context, name string, list model.List) error {
	c.mapMutex.Lock()
	defer c.mapMutex.Unlock()

	if c.lists == nil {
		return errNoData
	}

	if _, err := c.GetList(ctx, name); err != nil {
		return fmt.Errorf("list does not exist")
	} else if len(list.Words) == 0 {
		return fmt.Errorf("list must have at least one word")
	} else if len(list.Name) < 3 {
		return fmt.Errorf("list name is too short")
	}

	c.lists[listCollectionPrefix+name] = list
	return nil
}
