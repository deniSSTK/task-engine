package authService

import "errors"

var (
	FailedToValidateCredentials = errors.New("failed to validate credentials")
	InvalidCredentials          = errors.New("invalid credentials")
	EmailAlreadyExists          = errors.New("email already exists")

	FailedToCreateUser      = errors.New("failed to create user")
	FailedToUpdateUserInfo  = errors.New("failed to update user info")
	FailedToLogin           = errors.New("failed to login")
	FailedToGenerateTokens  = errors.New("failed to generate tokens")
	FailedToGetTokenPayload = errors.New("failed to get token payload")

	FailedToSaveUserSessionCache         = errors.New("failed to save user session in cache")
	FailedToVerifySessionExistingInCache = errors.New("failed to verify session existing in cache")

	FailedToAuthenticateUser = errors.New("failed to authenticate user")

	UndefinedUserStatus = errors.New("undefined user status")
)
