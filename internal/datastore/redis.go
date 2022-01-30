package datastore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"guess_my_word/internal/model"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient represents a Redis based Datastore client
type RedisClient struct {
	client *redis.Client
}

// NewRedis instantiates a new client
func NewRedis(addr string) *RedisClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client := redis.NewClient(&redis.Options{Addr: addr})
	status := client.Ping(ctx)
	if status.Err() != nil {
		log.Printf("Failed to create client: %v", status.Err())
		log.Println("Falling back to local mode. Generated data will not be persisted")
		client = nil
	}

	return &RedisClient{client: client}
}

// GetWord will retrieve a word for the given key
func (c *RedisClient) GetWord(ctx context.Context, key string) (model.Word, error) {
	if c.client == nil {
		return model.Word{}, errors.New("running in local mode. Data has not been retrieved")
	}

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
	if c.client == nil {
		return fmt.Errorf("running in local mode. Data has not been stored")
	}

	set, err := json.Marshal(word)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, wordCollectionPrefix+key, set, time.Hour*24*7).Err()
}
