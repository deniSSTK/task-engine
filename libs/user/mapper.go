package userDomain

import (
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	"github.com/google/uuid"
)

func MapRoleFromDomain(role UserRole) (authv1.UserRole, bool) {
	switch role {
	case RoleUser:
		return authv1.UserRole_USER_ROLE_USER, true
	case RoleAdmin:
		return authv1.UserRole_USER_ROLE_ADMIN, true
	default:
		return authv1.UserRole_USER_ROLE_UNSPECIFIED, false
	}
}

func MapRoleFromProto(role authv1.UserRole) (UserRole, bool) {
	switch role {
	case authv1.UserRole_USER_ROLE_USER:
		return RoleUser, true
	case authv1.UserRole_USER_ROLE_ADMIN:
		return RoleAdmin, true
	default:
		return RoleUser, false
	}
}

func MapProtoStatusFromDomain(status UserStatus) (authv1.UserStatus, bool) {
	switch status {
	case Active:
		return authv1.UserStatus_USER_STATUS_ACTIVE, true
	case Blocked:
		return authv1.UserStatus_USER_STATUS_BLOCKED, true
	default:
		return authv1.UserStatus_USER_STATUS_UNSPECIFIED, false
	}
}

func MapAuthUserFromProtoAuthUser(rawUser *authv1.AuthUser) (*AuthUser, bool) {
	id, err := uuid.Parse(rawUser.Id)
	if err != nil {
		return nil, false
	}

	role, ok := MapRoleFromProto(rawUser.Role)
	if !ok {
		return nil, false
	}

	return &AuthUser{
		Id:   id,
		Role: role,
	}, true
}
