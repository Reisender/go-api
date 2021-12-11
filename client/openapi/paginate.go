package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"go-api"
	"net/http"
)

var ErrStopPagination = fmt.Errorf("stop pagination")

func Paginate(ctx context.Context, c api.Client, endpoint string, params Params, page func([]byte) error) error {
	// parse the path
	path, err := c.NewURL(params.Prefix() + endpoint)
	if err != nil {
		return err
	}

	// add the params
	path.RawQuery = params.Values().Encode()

	// setup values for the pagination loop
	next := path.String()
	ok := true

	// pagination loop
	for ok && err == nil {
		// do the request
		var res *http.Response
		res, err = c.Get(ctx, next)
		if err != nil {
			return err
		}

		// parse the response
		resp := &Response{}
		dec := json.NewDecoder(res.Body)
		err = dec.Decode(resp)
		if err != nil {
			return err
		}
		err = res.Body.Close()
		if err != nil {
			return err
		}

		// handle the page results
		err = page(resp.Data)

		// see if there is a next page
		next, ok = resp.Links.Next()
	}

	if err == ErrStopPagination {
		return nil // don't pass this error on
	}

	return err
}
