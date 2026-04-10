package authGrpc

import "errors"

var (
	MissingAuthToken = errors.New("missing auth token")
)
