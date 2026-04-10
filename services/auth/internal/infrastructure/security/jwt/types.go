package jwt

import (
	userDomain "libs/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	RefreshExpiredAt time.Time
}

type TokenPayload struct {
	UserId uuid.UUID
	Role   userDomain.UserRole
}

type TokenClaims struct {
	Type TokenType
	jwt.RegisteredClaims
	TokenPayload
}
