package middleware_test

import (
	"errors"
	"go-api/middleware"
	"net/http"
	"testing"
)

const retries = 3

func TestRetryNilResponse(t *testing.T) {
	m := middleware.RetryOnStatusCodes(retries, middleware.StatusCodeRange{Low: 500, High: 599})
	tryCount := 0
	m(func(req *http.Request) (*http.Response, error) {
		tryCount++
		return nil, errors.New("TestRetryNilResponse error")
	})(nil)
	if tryCount != 2 {
		t.Errorf("expected %d retries, instead got %d", 2, tryCount)
	}
}

func TestRetry(t *testing.T) {
	m := middleware.RetryOnStatusCodes(retries, middleware.StatusCodeRange{Low: 500, High: 599})
	tryCount := 0
	m(func(req *http.Request) (*http.Response, error) {
		tryCount++
		return &http.Response{
			StatusCode: 500,
		}, nil
	})(nil)

	if tryCount != retries+1 {
		t.Errorf("expected %d retries, instead got %d", retries+1, tryCount)
	}
}

func TestNoRetry(t *testing.T) {
	m := middleware.RetryOnStatusCodes(retries, middleware.StatusCodeRange{Low: 500, High: 599})
	tryCount := 0
	m(func(req *http.Request) (*http.Response, error) {
		tryCount++
		return &http.Response{
			StatusCode: 200,
		}, nil
	})(nil)
	if tryCount != 1 {
		t.Errorf("expected %d retry, instead got %d", 1, tryCount)
	}
}
