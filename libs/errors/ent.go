package defErrors

import "errors"

var (
	FailedToCreateEntClient = "failed to create ent client"

	FailedToGetData = errors.New("failed to update data")
)
