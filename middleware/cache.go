package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Reisender/go-api"
)

type Mode int

const (
	ModeDefault = iota
	ModeCacheOnly
)

type GetSetter interface {
	Get(key string) (io.ReadCloser, error)
	Set(key string, ttl time.Duration, body io.ReadCloser) error
}

// Header will set the header name and value
func Cache(mode Mode, ttl time.Duration, store GetSetter) api.Middleware {
	// create a middleware func
	return func(next api.Dofn) api.Dofn {

		// return a new Dofn
		return func(req *http.Request) (*http.Response, error) {
			// TOD: do the cache logic

			key := getCacheKey(req)
			headersKey := fmt.Sprintf("%s-headers", key)

			if mode == ModeCacheOnly {
				// check for the cache files
				headersReader, err := store.Get(headersKey)
				if err != nil {
					return &http.Response{
						StatusCode: 404,
						Body:       io.NopCloser(bytes.NewReader([]byte("Not found"))),
					}, nil
				}
				body, err := store.Get(key)
				if err != nil {
					return &http.Response{
						StatusCode: 404,
						Body:       io.NopCloser(bytes.NewReader([]byte("Not found"))),
					}, nil
				}

				headers := make(http.Header)
				headersJson, err := io.ReadAll(headersReader)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(headersJson, &headers)
				if err != nil {
					return nil, err
				}

				return &http.Response{
					Request:    req,
					Body:       body,
					Header:     headers,
					StatusCode: 200,
					Status:     "200 OK",
					Proto:      req.Proto,
					ProtoMajor: req.ProtoMajor,
					ProtoMinor: req.ProtoMinor,
				}, nil
			}

			resp, err := next(req)
			if err != nil {
				return nil, err
			}

			headers, err := json.Marshal(resp.Header)
			if err == nil {
				// Read the response body
				bodyBytes, err := io.ReadAll(resp.Body)
				if err == nil {
					// Close the original body
					resp.Body.Close()

					// Create two new copies of the body
					resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					cacheBody := io.NopCloser(bytes.NewBuffer(bodyBytes))

					// Store in cache
					store.Set(key, ttl, cacheBody)
					store.Set(headersKey, ttl, io.NopCloser(bytes.NewReader(headers)))
				}
			}

			return resp, nil
		}

	}
}

func getCacheKey(req *http.Request) string {
	if req == nil {
		return ""
	}

	// Get full URL
	urlStr := req.URL.String()

	h := md5.New()
	io.WriteString(h, urlStr)

	// If body exists and is readable, include it in the hash
	if req.Body != nil {
		// Clone the body so we don't consume it
		bodyBytes, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err == nil {
			// Put the body back for future readers
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// Add body to hash
			h.Write(bodyBytes)
		}
	}

	return hex.EncodeToString(h.Sum(nil))
}
