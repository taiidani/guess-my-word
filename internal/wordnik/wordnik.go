package wordnik

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const apiKeyTODO = "" // TODO: Remove

var urlPrefix = "https://api.wordnik.com/v4" // The prefix added to all requests

type (
	service interface {
	}

	// Client represents a client for the Wordnik API
	Client struct {
		apiKey string // API key for authenticating to Wordnik
		client *http.Client
		Words  Words
	}

	genericResponse struct {
		Message string
	}
)

// NewClient returns a new instance of the Wordnik API client
func NewClient(apiKey string) *Client {
	client := &Client{
		apiKey: apiKey,
		client: http.DefaultClient,
	}
	client.Words = client.newWords()

	return client
}

func (c *Client) get(ctx context.Context, url string, params url.Values) (*http.Request, error) {
	// Set up the URL
	url = urlPrefix + "/" + strings.TrimPrefix(url, "/")

	// Set up the parameters
	params.Add("api_key", c.apiKey)
	url = url + "?" + params.Encode()

	log.Println(url)

	r, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	return r, err
}

func (c *Client) handleRequest(req *http.Request, target interface{}) error {
	response, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		output := genericResponse{}
		if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
			output.Message = "Unable to decode response message from Wordnik API"
		}

		return fmt.Errorf("%d: %s", response.StatusCode, output.Message)
	}

	return json.NewDecoder(response.Body).Decode(&target)
}
