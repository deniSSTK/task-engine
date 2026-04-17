package defErrors

import "errors"

var (
	UserUnauthenticated = errors.New("user is not authenticated")

	MissingMetadata      = errors.New("missing metadata")
	MissingMetadataValue = errors.New("missing metadata value")

	MetadataOptionUnspecified = errors.New("metadata option unspecified")

	MethodOptionsNotFound = errors.New("method options not found")
	MissingRequiredOption = errors.New("missing required option")
	InvalidExtensionType  = errors.New("invalid extension type")
)
