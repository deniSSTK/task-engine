package authGrpc

import (
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/security/jwt"
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
)

func MapTokenPairToProtoTokenDetail(data *jwt.TokenPair) *authv1.TokenDetails {
	return &authv1.TokenDetails{
		AccessToken:      data.AccessToken,
		RefreshToken:     data.RefreshToken,
		RefreshExpiresAt: data.RefreshExpiredAt.Unix(),
	}
}
