package datastore

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

const (
	defaultProjectID     = "burnished-fold-276716"
	wordCollectionPrefix = "words/"
)

// Client represents a Datastore client
type Client struct{}

var (
	// LocalOnly will not perform any datastore operations if true
	// It will also not return any errors for Set operations
	LocalOnly bool

	// The local GCP Firestore client
	client *firestore.Client
)

func init() {
	var err error
	ctx := context.Background()

	// Default the authentication location to an auth.json file in the local directory
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./auth.json")
	}

	// Default the project ID, allowing overrides during development
	if os.Getenv("GOOGLE_APPLICATION_PROJECT_ID") == "" {
		os.Setenv("GOOGLE_APPLICATION_PROJECT_ID", defaultProjectID)
	}

	client, err = firestore.NewClient(ctx, os.Getenv("GOOGLE_APPLICATION_PROJECT_ID"))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		log.Println("Falling back to local mode. Generated data will not be persisted")
		LocalOnly = true
	}
}

// New instantiates a new Datastore client
func New() *Client {
	return &Client{}
}

// GetWord will retrieve a word for the given key
func (c *Client) GetWord(ctx context.Context, key string, word interface{}) (err error) {
	if LocalOnly {
		return errors.New("Running in local mode. Data has not been retrieved")
	}

	doc := client.Doc(wordCollectionPrefix + key)
	data, err := doc.Get(ctx)
	if err != nil {
		return err
	}

	err = data.DataTo(&word)
	return
}

// SetWord will store a Word for the given key
func (c *Client) SetWord(ctx context.Context, key string, word interface{}) error {
	if LocalOnly {
		log.Println("Running in local mode. Data has not been stored")
		return nil
	}

	_, err := client.Doc(wordCollectionPrefix+key).Create(ctx, word)
	return err
}
