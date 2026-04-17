package grpcAuth

import "errors"

var (
	MissingAuthToken = errors.New("missing auth token")

	MissingAuthUser = errors.New("missing auth user")
)
