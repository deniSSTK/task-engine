package defErrors

import "errors"

var (
	FailedMap = errors.New("failed to map")

	InvalidCredentials = errors.New("invalid credentials")
)
