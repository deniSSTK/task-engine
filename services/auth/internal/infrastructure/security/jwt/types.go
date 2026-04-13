package jwt

import (
	"time"

	userDomain "github.com/deniSSTK/task-engine/libs/user"
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
