package authGrpc

import (
	"auth-service/internal/infra/security/jwt"
	proto "proto/auth"
)

func MapTokenPairToProtoTokenResponse(data *jwt.TokenPair) *proto.TokensResponse {
	return &proto.TokensResponse{
		AccessToken:      data.AccessToken,
		RefreshToken:     data.RefreshToken,
		RefreshExpiresAt: data.RefreshExpiredAt.Unix(),
	}
}
