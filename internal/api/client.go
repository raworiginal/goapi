// Package api is the abstraction for making http calls.
package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

type Response struct {
	StatusCode int
	Body       []byte
	Duration   time.Duration
}

type Config struct {
	Timeout time.Duration
}

type Client interface {
	Do(method, url string, body []byte) (*Response, error)
}

type HTTPClient struct {
	config Config
}

func NewHTTPClient(config Config) Client {
	return &HTTPClient{config: config}
}

func (c *HTTPClient) Do(method, url string, body []byte) (*Response, error) {
	// Create a context with Timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	// Create the HTTP Request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// start timing
	start := time.Now()

	// send the Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Calculate Duration
	duration := time.Since(start)

	// Read Response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Return the respose
	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Duration:   duration,
	}, nil
}
