package openapi

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/Reisender/go-api"
)

func Count(ctx context.Context, c api.Client, endpoint string, params Params) (int, error) {
	path := ""

	if !isNil(params) {
		path += params.Prefix()
	}
	path += endpoint

	reqURL, err := c.NewURL(path)
	if err != nil {
		return 0, err
	}

	vals := url.Values{}
	if !isNil(params) {
		vals = params.Values()
	}
	vals.Del("limit")         // limit isn't allowed for this kind of request
	vals.Set("count", "true") // force the count param
	reqURL.RawQuery = vals.Encode()

	res, err := c.Get(ctx, reqURL.String())
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	countResponse := struct {
		Count int `json:"count"`
	}{}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&countResponse)
	if err != nil {
		return 0, err
	}

	return countResponse.Count, nil
}
