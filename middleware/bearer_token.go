package middleware

import (
	"net/http"

	"github.com/Reisender/go-api"
)

func BearerToken(token string) api.Middleware {
	// create a middleware func
	return func(next api.Dofn) api.Dofn {

		// return a new Dofn
		return func(req *http.Request) (*http.Response, error) {
			req.Header.Add("Authorization", "Bearer "+token)
			return next(req)
		}

	}
}
