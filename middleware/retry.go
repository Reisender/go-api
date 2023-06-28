package middleware

import (
	"fmt"
	"go-api"
	"net/http"
	"time"
)

// ErrMaxRetries is the error that represents when the max number of retries has been reached
type ErrMaxRetries struct {
	Err error // the wrapped error
}

// Error implements the error interface
func (e ErrMaxRetries) Error() string {
	return fmt.Sprintf("max retries reached: %s", e.Err)
}
func (e ErrMaxRetries) Unwrap() error {
	return e.Err
}

// RetryOnStatusCodes is a Do func middleware that will retry based on status codes
func RetryOnStatusCodes(retry uint, statusCodes ...StatusCodeRange) api.Middleware {
	// return the middleware func
	return func(next api.Dofn) api.Dofn {

		// return the Do func
		return func(req *http.Request) (*http.Response, error) {

			// If there is an error, resp can be nil
			resp, err := next(req)

			retryCount := uint(0)
			for retryCount < retry && (resp == nil || InRanges(resp.StatusCode, statusCodes)) {
				retryCount++
				resp, err = next(req)
				if err != nil {
					return nil, ErrMaxRetries{err}
				}
			}

			return resp, err
		}

	}
}

// Non2XXStatusCodes is the status code range representing non 2XX status codes
var Non2XXStatusCodes = StatusCodeRange{Low: 300, High: 599}

// RetryWithDelay is a Do func middleware that will retry based on status codes or on err.
// It also can take a delay min, max, and backoff multiplier. If no range is passed, it defaults
// to Non2XXStatusCodes range
func RetryWithDelay(retry uint, delayMin, delayMax time.Duration, delayRamp float32, ranges ...StatusCodeRange) api.Middleware {
	// default to the Non2XXStatusCodes
	if len(ranges) == 0 {
		ranges = []StatusCodeRange{Non2XXStatusCodes}
	}

	// return the middleware func
	return func(next api.Dofn) api.Dofn {

		// return the Do func
		return func(req *http.Request) (*http.Response, error) {

			// If there is an error, resp can be nil
			resp, err := next(req)

			// retry on status code >= 300 or err from next
			retryCount := uint(0)
			delay := delayMin
			for retryCount < retry && (err != nil || resp == nil || InRanges(resp.StatusCode, ranges)) {
				retryCount++
				select {
				case <-req.Context().Done():
					return nil, req.Context().Err()
				case <-time.After(delay):
					// update the delay for the backoff
					if delay < delayMax {
						delay = time.Duration(float32(delay) * delayRamp)
					}
					if delay > delayMax {
						delay = delayMax
					}

					// try again
					resp, err = next(req)
				}
			}

			if resp != nil && InRanges(resp.StatusCode, ranges) {
				err = ErrMaxRetries{ErrStatusCode{resp.Status, resp.StatusCode}}
			} else if err != nil {
				err = ErrMaxRetries{err}
			}

			return resp, err
		}

	}
}
