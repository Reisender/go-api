package middleware

import (
	"go-api"
	"net/http"
)

type StatusCodeRange struct {
	Low  int
	High int
}

func (scr StatusCodeRange) InRange(code int) bool {
	return (code >= scr.Low && code <= scr.High)
}

func checkStatusCodes(code int, checkCodes []StatusCodeRange) bool {
	for _, codeRange := range checkCodes {
		if codeRange.InRange(code) {
			return true
		}
	}

	return false
}

// RetryOnStatusRange is a Do func middleware that will retry based on status codes
func RetryOnStatusCodes(retry uint, statusCodes ...StatusCodeRange) api.Middleware {
	// return the middleware func
	return func(next api.Dofn) api.Dofn {

		// return the Do func
		return func(req *http.Request) (*http.Response, error) {

			resp, err := next(req)

			retryCount := uint(0)
			for retryCount < retry && checkStatusCodes(resp.StatusCode, statusCodes) {
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
