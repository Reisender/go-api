package middleware

import (
	"net/http"

	"github.com/Reisender/go-api"
)

// StatusCodeRange is a range of status codes
type StatusCodeRange struct {
	Low  int
	High int
}

// InRange checks a status code to see if it is in the range
func (scr StatusCodeRange) InRange(code int) bool {
	return (code >= scr.Low && code <= scr.High)
}

// InRanges checks to see if the status code is in the ranges
func InRanges(code int, ranges []StatusCodeRange) bool {
	for _, codeRange := range ranges {
		if codeRange.InRange(code) {
			return true
		}
	}

	return false
}

// ErrStatusCode This is used when you want to handle a certain status code
// as an error.
type ErrStatusCode struct {
	Status string
	Code   int
}

func (esc ErrStatusCode) Error() string {
	return esc.Status
}

// ErrorOnStatusCodes will return an error on certain status codes
func ErrorOnStatusCodes(statusCodes ...StatusCodeRange) api.Middleware {
	// return the middleware func
	return func(next api.Dofn) api.Dofn {

		// return the Do func
		return func(req *http.Request) (*http.Response, error) {
			// pass the request on to the next
			res, err := next(req)

			// make sure there isn't already an error
			if err != nil {
				return res, err
			}

			// now see if it is an error code to convert to error
			if InRanges(res.StatusCode, statusCodes) {
				return res, ErrStatusCode{res.Status, res.StatusCode}
			}

			return res, err
		}

	}
}
