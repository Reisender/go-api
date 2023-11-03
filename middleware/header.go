package middleware

import (
	"net/http"

	"github.com/Reisender/go-api"
)

// Header will set the header name and value
func Header(name, value string) api.Middleware {
	// create a middleware func
	return func(next api.Dofn) api.Dofn {

		// return a new Dofn
		return func(req *http.Request) (*http.Response, error) {
			req.Header.Add(name, value)
			return next(req)
		}

	}
}
