package middleware_test

import (
	"context"
	"errors"
	"go-api/middleware"
	"net/http"
	"testing"
	"time"
)

const retries = 3

func TestRetryOnStatusCodes(t *testing.T) {
	t.Run("retry err response", func(t *testing.T) {
		m := middleware.RetryOnStatusCodes(retries, middleware.StatusCodeRange{Low: 500, High: 599})
		tryCount := 0
		m(func(req *http.Request) (*http.Response, error) {
			tryCount++
			return nil, errors.New("TestRetryNilResponse error")
		})(nil)
		if tryCount != 2 {
			t.Errorf("expected %d retries, instead got %d", 2, tryCount)
		}
	})

	t.Run("retry", func(t *testing.T) {
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
	})

	t.Run("no retry", func(t *testing.T) {
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
	})
}

func TestRetryWithDelay(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodHead, "", nil)
	retries := uint(10)

	delayMin := time.Millisecond
	delayMax := time.Millisecond * 2
	delayRamp := float32(2.0)
	m := middleware.RetryWithDelay(retries, delayMin, delayMax, delayRamp)

	tryCount := uint(0)
	var start time.Time

	m(func(req *http.Request) (*http.Response, error) {
		got := time.Now().Sub(start)
		want := time.Duration(float32(tryCount-1) * delayRamp * float32(delayMin))
		if want > delayMax {
			want = delayMax // don't want anything above max
		}
		start = time.Now() // start the clock on the next run
		if tryCount == 0 {
			// don't check things on the first run
		} else if want > got {
			t.Errorf("try %d expected a delay of at least %v but got %v", tryCount, want, got)
		} else if got > delayMax*2 { // *2 to give the max some breathing room. Even with processing time it shouldn't be more that double
			t.Errorf("try %d expected a delay max around %v but got %v", tryCount, delayMax, got)
		}
		tryCount++
		return &http.Response{
			StatusCode: 500,
		}, nil
	})(req)

	if tryCount != retries+1 {
		t.Errorf("expected %d retries, instead got %d", retries+1, tryCount)
	}
}
