package jwt

import "errors"

var (
	FailedToGenerateToken = errors.New("failed to generate token")

	InvalidToken     = errors.New("invalid token")
	InvalidTokenType = errors.New("invalid token type")
)
