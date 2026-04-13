package grpcAuth

import "errors"

var (
	AdminOnly = errors.New("admin only")

	MissingAuthToken = errors.New("missing auth token")

	MissingAuthUser = errors.New("missing auth user")
)
