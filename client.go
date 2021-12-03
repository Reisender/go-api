package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Dofn is the func def for the http.Client{}.Do func
type Dofn func(req *http.Request) (*http.Response, error)

type Middleware func(doer Dofn) Dofn

type Client struct {
	httpClient *http.Client
	host       string
	do         Dofn
}

// NewClient creates a new Client
// The (optional) Middleware Do funcs will run in sequence
// when the Do func is called with a request.
func NewClient(baseEndpoint string, timeout time.Duration, doers ...Middleware) *Client {
	client := &http.Client{
		Timeout: timeout,
	}

	c := &Client{
		httpClient: client,
		host:       baseEndpoint,
	}

	//
	// construct the Do func
	//

	// start with the base Do func
	do := c.httpClient.Do

	// apply the middleware Do funcs
	// in reverse order so that then end up executing
	// in the ordered they were passed in
	for i := len(doers) - 1; i >= 0; i-- {
		do = doers[i](do)
	}

	// use the wrapped Do func for the client
	c.do = do

	return c
}

// NewRequestWithContext wraps the http version and sets the url.
// This allows the endpoint being passed to not include the host or base part of the url.
func (c Client) NewRequestWithContext(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf(
		"%s/%s",
		strings.TrimRight(c.host, "/"),
		strings.TrimLeft(endpoint, "/"),
	)

	return http.NewRequestWithContext(ctx, method, url, body)
}

// Get is similar to http.Client{}.Get except it uses the BearerTokenClient
// and defaults to JSON as the payload to and from the server.
func (c Client) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	req, err := c.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Do calles the Do func that was built up from the middleware
func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}
