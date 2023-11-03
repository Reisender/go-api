package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"

	"github.com/Reisender/go-api"
)

// NewMock mocks the Do func and doesn't
// call in to the next Dofn in line.
func NewMock(mock api.Dofn) api.Middleware {
	// return the middleware func
	return func(next api.Dofn) api.Dofn {

		// return the Do func
		return func(req *http.Request) (*http.Response, error) {

			// only do our mock Dofn
			return mock(req)
		}

	}
}

// MockResponse mocks a response object and can be used with NewMock.
func MockResponse(mock func(req *http.Request) (statusCode int, respBody string)) api.Dofn {
	// handler that returns the resonse directly and doesn't pass on to the next
	return func(req *http.Request) (*http.Response, error) {
		statusCode, respBody := mock(req)
		buf := bytes.NewBufferString(fmt.Sprintf("HTTP/1.1 %d\n\n%s", statusCode, respBody))
		return http.ReadResponse(bufio.NewReader(buf), req)
	}
}

// NewMockResponse is a convenience func that combines NewMock and MockResponse
func NewMockResponse(mock func(req *http.Request) (statusCode int, respBody string)) api.Middleware {
	// return the middleware func
	return NewMock(MockResponse(mock))
}
