package httpclient

import (
	"io"
	"net/http"
	"time"
)

// Client is a wrapper around http.Client
type Client struct {
	httpClient *http.Client
}

// NewHttpClient creates a new instance of Client with a specified timeout
func NewHttpClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get performs a GET request
func (c *Client) Get(url string) (*http.Response, error) {
	return c.httpClient.Get(url)
}

// Post performs a POST request
func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return c.httpClient.Post(url, contentType, body)
}
