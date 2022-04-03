package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"guess_my_word/internal/model"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

// RedisClient represents a Redis based Datastore client
type RedisClient struct {
	client redisClient
	mutex  sync.Mutex
}

// NewRedis instantiates a new client
func NewRedis(addr string) *RedisClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client := redis.NewClient(&redis.Options{Addr: addr})
	status := client.Ping(ctx)
	if status.Err() != nil {
		log.Printf("Failed to create Redis client")
		panic(status.Err())
	}

	return &RedisClient{client: client, mutex: sync.Mutex{}}
}

// GetWord will retrieve a word for the given key
func (c *RedisClient) GetWord(ctx context.Context, key string) (model.Word, error) {
	ret := model.Word{}
	doc := c.client.Get(ctx, wordCollectionPrefix+key)
	if doc.Err() != nil {
		return ret, doc.Err()
	}

	err := json.Unmarshal([]byte(doc.Val()), &ret)
	return ret, err
}

// SetWord will store a Word for the given key
func (c *RedisClient) SetWord(ctx context.Context, key string, word model.Word) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	set, err := json.Marshal(word)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, wordCollectionPrefix+key, set, time.Hour*24*7).Err()
}

func (c *RedisClient) GetLists(ctx context.Context) ([]string, error) {
	keys := c.client.Keys(ctx, listCollectionPrefix+"*")
	if keys.Err() != nil {
		return nil, keys.Err()
	}

	ret := keys.Val()
	sort.Strings(ret)
	for i := range ret {
		ret[i] = strings.TrimPrefix(ret[i], listCollectionPrefix)
	}

	return ret, nil
}

func (c *RedisClient) GetList(ctx context.Context, name string) (model.List, error) {
	ret := model.List{}
	list := c.client.Get(ctx, listCollectionPrefix+name)
	if list.Err() != nil {
		return ret, list.Err()
	}

	err := json.Unmarshal([]byte(list.Val()), &ret)
	return ret, err
}

func (c *RedisClient) CreateList(ctx context.Context, name string, list model.List) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if exist := c.client.Exists(ctx, listCollectionPrefix+name); exist.Err() != nil {
		return fmt.Errorf("could not check for existing list: %w", exist.Err())
	} else if exist.Val() > 0 {
		return fmt.Errorf("list already exists for that name")
	}

	set, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, listCollectionPrefix+name, set, 0).Err()
}

func (c *RedisClient) DeleteList(ctx context.Context, name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.client.Del(ctx, name).Err()
}

func (c *RedisClient) UpdateList(ctx context.Context, name string, list model.List) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	set, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, listCollectionPrefix+name, set, 0).Err()
}
