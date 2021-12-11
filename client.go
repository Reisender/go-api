package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Dofn is the func def for the http.Client{}.Do func
type Dofn func(req *http.Request) (*http.Response, error)

type Middleware func(doer Dofn) Dofn

type Client interface {
	NewURL(endpoint string) (*url.URL, error)
	NewRequestWithContext(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error)
	Get(ctx context.Context, endpoint string) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

type BaseClient struct {
	httpClient *http.Client
	host       string
	base       string
	do         Dofn
}

// NewClient creates a new BaseClient
// The (optional) Middleware Do funcs will run in sequence
// when the Do func is called with a request.
// The baseEndpoint is not added automatically with Do or Get etc...
// It is added in if you use the BaseClient.NewURL to generate your new URL.
// This allows the client to be used with or without assuming the base part.
// This is useful if you are using the "links" part of responses which already
// have the base part in them.
func NewClient(host, baseEndpoint string, timeout time.Duration, doers ...Middleware) *BaseClient {
	client := &http.Client{
		Timeout: timeout,
	}

	c := &BaseClient{
		httpClient: client,
		host:       host,
		base:       baseEndpoint,
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

func (c BaseClient) NewURL(endpoint string) (*url.URL, error) {
	return url.Parse(c.base + endpoint)
}

// NewRequestWithContext wraps the http version and sets the url.
// This allows the endpoint being passed to not include the host or base part of the url.
func (c BaseClient) NewRequestWithContext(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf(
		"%s/%s",
		strings.TrimRight(c.host, "/"),
		strings.TrimLeft(endpoint, "/"),
	)

	return http.NewRequestWithContext(ctx, method, url, body)
}

// Get is similar to http.Client{}.Get except it uses the BearerTokenClient
// and defaults to JSON as the payload to and from the server.
func (c BaseClient) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	req, err := c.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Do calles the Do func that was built up from the middleware
func (c BaseClient) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}
