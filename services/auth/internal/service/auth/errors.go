package authService

import "errors"

var (
	FailedToValidateCredentials = errors.New("failed to validate credentials")
	EmailAlreadyExists          = errors.New("email already exists")
	FailedToCreateUser          = errors.New("failed to create user")

	FailedToSaveUserSessionCache = errors.New("failed to save user session in cache")
)
