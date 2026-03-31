package authGrpc

import (
	"auth-service/internal/infra/security/jwt"
	proto "proto/proto/auth/v1"
)

func MapTokenPairToProtoTokenDetail(data *jwt.TokenPair) *proto.TokenDetails {
	return &proto.TokenDetails{
		AccessToken:      data.AccessToken,
		RefreshToken:     data.RefreshToken,
		RefreshExpiresAt: data.RefreshExpiredAt.Unix(),
	}
}
