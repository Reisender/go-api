package middleware

import (
	"go-api"
	"net/http"
)

// RetryOnStatusRange is a Do func middleware that will retry based on status codes
func RetryOnStatusCodes(retry uint, statusCodes ...StatusCodeRange) api.Middleware {
	// return the middleware func
	return func(next api.Dofn) api.Dofn {

		// return the Do func
		return func(req *http.Request) (*http.Response, error) {

			resp, err := next(req)

			retryCount := uint(0)
			for retryCount < retry && InRanges(resp.StatusCode, statusCodes) {
				retryCount++
				resp, err = next(req)
				if err != nil {
					return nil, err
				}
			}

			return resp, err
		}

	}
}
