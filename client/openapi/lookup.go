package openapi

import (
	"context"
	"encoding/json"

	"github.com/Reisender/go-api"
)

// Lookup is a common implementation of a single resource lookup.
// This typically looks like /users/74 that returns a single resource
func Lookup(ctx context.Context, c api.Client, endpoint string, resource interface{}) error {
	url, err := c.NewURL(endpoint)
	if err != nil {
		return err
	}

	res, err := c.Get(ctx, url.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resp := &Response{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp.Data, resource)
}
