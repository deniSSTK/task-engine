package authService

import "errors"

var (
	FailedToValidateCredentials = errors.New("failed to validate credentials")
	InvalidCredentials          = errors.New("invalid credentials")
	EmailAlreadyExists          = errors.New("email already exists")

	FailedToCreateUser     = errors.New("failed to create user")
	FailedToUpdateUserInfo = errors.New("failed to update user info")
	FailedToLogin          = errors.New("failed to login")
	FailedToGenerateTokens = errors.New("failed to generate tokens")

	FailedToSaveUserSessionCache = errors.New("failed to save user session in cache")
)
