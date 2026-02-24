package api

import (
	"strings"
	"testing"
	"time"
)

func TestGetRequest(t *testing.T) {
	// Arrange: Set up test data
	config := Config{Timeout: 5 * time.Second}
	client := NewHTTPClient(config)

	// Act: Do Something
	response, err := client.Do("GET", "https://httpbin.org/get", nil)
	// Assert: Check results
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", response.StatusCode)
	}

	if len(response.Body) == 0 {
		t.Error("expected non-empty body")
	}
}

func TestPostRequest(t *testing.T) {
	config := Config{Timeout: 5 * time.Second}
	client := NewHTTPClient(config)

	response, err := client.Do("POST", "https://httpbin.org/post", []byte("this is a test"))
	if err != nil {
		t.Fatalf("expected no error, got %v", response.StatusCode)
	}

	if response.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", err)
	}

	if !strings.Contains(string(response.Body), "this is a test") {
		t.Error("Body does not contain test data")
	}
}

func TestTimeout(t *testing.T) {
	config := Config{Timeout: 2 * time.Second}
	client := NewHTTPClient(config)

	response, err := client.Do("GET", "https://httpbin.org/delay/5", nil)
	// Assert: Should have an error (timeout)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	// Check it's a context deadline error
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("expected context deadline error, got: %v", err)
	}

	// Response should be nil since it timed out
	if response != nil {
		t.Error("expected nil response on timeout")
	}
}
