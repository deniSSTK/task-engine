package authGrpcV1

import (
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/security/jwt"
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
)

func MapTokenPairToProtoTokenDetail(data *jwt.TokenPair) *authv1.TokenDetails {
	return &authv1.TokenDetails{
		AccessToken:      data.AccessToken,
		RefreshToken:     data.RefreshToken,
		RefreshExpiresAt: data.RefreshExpiredAt.Unix(),
	}
}

func MapDomainUserToProtoUser(data *userDomain.User) (*authv1.User, bool) {
	role, ok := userDomain.MapRoleFromDomain(data.Role)
	if !ok {
		return nil, false
	}

	status, ok := userDomain.MapProtoStatusFromDomain(data.Status)
	if !ok {
		return nil, false
	}

	return &authv1.User{
		Id:         data.Id.String(),
		Name:       data.Name,
		SecondName: data.SecondName,
		FullName:   data.FullName,
		Email:      data.Email,
		Role:       role,
		Status:     status,
	}, true
}
