package datastore

import (
	"context"
	"errors"
	"fmt"
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
func (c *RedisClient) GetWord(ctx context.Context, key string) (string, error) {
	if c.client == nil {
		return "", errors.New("running in local mode. Data has not been retrieved")
	}

	doc := c.client.Get(ctx, wordCollectionPrefix+key)
	return doc.Val(), doc.Err()
}

// SetWord will store a Word for the given key
func (c *RedisClient) SetWord(ctx context.Context, key string, word string) error {
	if c.client == nil {
		return fmt.Errorf("Running in local mode. Data has not been stored")
	}

	return c.client.Set(ctx, wordCollectionPrefix+key, word, time.Hour*24*7).Err()
}
