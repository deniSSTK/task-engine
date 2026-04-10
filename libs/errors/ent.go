package defErrors

import "errors"

var (
	FailedToCreateEntClient = "failed to create ent client"

	FailedToGetData = errors.New("failed to update data")

	MissingMetadata      = errors.New("missing metadata")
	MissingMetadataValue = errors.New("missing metadata value")
)
