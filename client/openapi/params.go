package openapi

import "net/url"

type Params interface {
	// Values converts the params to the url.Values type
	Values() url.Values

	// Prefix is to define a prefix to the endpoint
	Prefix() string
}
