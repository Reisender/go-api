package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

// MockCache implements the GetSetter interface for testing
type MockCache struct {
	cache map[string][]byte
}

func NewMockCache() *MockCache {
	return &MockCache{
		cache: make(map[string][]byte),
	}
}

func (m *MockCache) Get(key string) (io.ReadCloser, error) {
	if data, ok := m.cache[key]; ok {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	return nil, io.EOF
}

func (m *MockCache) Set(key string, ttl time.Duration, body io.ReadCloser) error {
	if body == nil {
		return nil
	}
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	m.cache[key] = data
	return nil
}

func TestCacheDefault(t *testing.T) {
	// Create a mock cache
	store := NewMockCache()

	// Create a mock handler that returns a predictable response
	handler := func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("test response")),
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
		}, nil
	}

	// Apply the cache middleware
	middleware := Cache(ModeDefault, 1*time.Minute, store)
	wrappedHandler := middleware(handler)

	// Create a test request
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// First request should go through to the handler and cache the result
	resp1, err := wrappedHandler(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp1.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp1.StatusCode)
	}

	// Check that response was stored in cache
	key := getCacheKey(req)
	headersKey := key + "-headers"
	if _, ok := store.cache[key]; !ok {
		t.Error("Response body not stored in cache")
	}
	if _, ok := store.cache[headersKey]; !ok {
		t.Error("Response headers not stored in cache")
	}

	// Read the response body and save it to a variable
	bodyBytes, err := io.ReadAll(resp1.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	bodyStr := string(bodyBytes)

	if bodyStr != "test response" {
		t.Errorf("Expected 'test response', got '%s'", bodyStr)
	}
}

func TestCacheCacheOnly(t *testing.T) {
	// Create a mock cache
	store := NewMockCache()

	// Pre-populate the cache
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	key := getCacheKey(req)
	headersKey := key + "-headers"

	// Store a response body
	store.cache[key] = []byte("cached response")

	// Store headers
	headers := http.Header{"Content-Type": []string{"text/plain"}}
	headerData, _ := json.Marshal(headers)
	store.cache[headersKey] = headerData

	// Create a handler that should never be called
	handler := func(req *http.Request) (*http.Response, error) {
		t.Fatal("Handler should not be called in CacheOnly mode")
		return nil, nil
	}

	// Apply the cache middleware in CacheOnly mode
	middleware := Cache(ModeCacheOnly, 1*time.Minute, store)
	wrappedHandler := middleware(handler)

	// Request should be served from cache
	resp, err := wrappedHandler(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Read the response body
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "cached response" {
		t.Errorf("Expected 'cached response', got '%s'", string(body))
	}

	// Check if Content-Type header was correctly restored
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got '%s'", contentType)
	}
}

func TestCacheCacheOnlyNotFound(t *testing.T) {
	// Create an empty mock cache
	store := NewMockCache()

	// Create a handler that should never be called
	handler := func(req *http.Request) (*http.Response, error) {
		t.Fatal("Handler should not be called in CacheOnly mode")
		return nil, nil
	}

	// Apply the cache middleware in CacheOnly mode
	middleware := Cache(ModeCacheOnly, 1*time.Minute, store)
	wrappedHandler := middleware(handler)

	// Create a test request
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// Request should return 404 since nothing is in cache
	resp, err := wrappedHandler(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}
}

func TestGetCacheKey(t *testing.T) {
	// Test with URL only
	req1, _ := http.NewRequest("GET", "http://example.com/path?q=test", nil)
	key1 := getCacheKey(req1)
	if key1 == "" {
		t.Error("Expected non-empty cache key")
	}

	// Test with URL and body
	body := bytes.NewBufferString("request body")
	req2, _ := http.NewRequest("POST", "http://example.com/path", body)
	key2 := getCacheKey(req2)
	if key2 == "" {
		t.Error("Expected non-empty cache key")
	}

	// Keys should be different
	if key1 == key2 {
		t.Error("Expected different keys for different requests")
	}

	// Test that same request generates same key
	req3, _ := http.NewRequest("GET", "http://example.com/path?q=test", nil)
	key3 := getCacheKey(req3)
	if key1 != key3 {
		t.Error("Expected same key for identical requests")
	}
}
