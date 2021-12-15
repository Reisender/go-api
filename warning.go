package api

import "errors"

// this can be users to indicate errors that are
// more like warning and should be logged, but
// should not stop execution flow.
var Warning = errors.New("warning")
