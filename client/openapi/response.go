package openapi

import "encoding/json"

type Response struct {
	Data  json.RawMessage `json:"data,omitempty"` // delay parsing here to be parsed by the specific api call
	Links Links           `json:"links"`
}
