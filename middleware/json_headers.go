package middleware

import (
	"go-api"
	"net/http"
)

/*
func NewExampleMiddleware(val string) api.Middleware {
	// return the middleware func
	return func(next api.Dofn) api.Dofn {

		// return the Do func
		return func(req *http.Request) (*http.Response, error) {

			// ... do the things this middleware does with val

			// pass the request on to the next
			return next(req)
		}

	}
}
*/

// JSONPayload is a Do func middleware that sets the JSON payload headers
func JSONHeaders(next api.Dofn) api.Dofn {
	return func(req *http.Request) (*http.Response, error) {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		return next(req)
	}
}
