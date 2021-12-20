package middleware

import (
	"go-api"
	"net/http"
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
