package middleware

import (
	"fmt"
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

func InRanges(code int, ranges []StatusCodeRange) bool {
	for _, codeRange := range ranges {
		if codeRange.InRange(code) {
			return true
		}
	}

	return false
}

// This is used when you want to handle a certain status code
// as an error.
type ErrStatusCode struct {
	Code int
}

func (esc ErrStatusCode) Error() string {
	switch esc.Code {
	case http.StatusBadRequest:
		return "bad request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not found"
	case http.StatusInternalServerError:
		return "internal server error"
	}

	return fmt.Sprintf("error status code %d", esc.Code)
}

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
				return res, ErrStatusCode{res.StatusCode}
			}

			return res, err
		}

	}
}
