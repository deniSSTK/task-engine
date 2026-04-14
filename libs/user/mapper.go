package userDomain

import authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"

func MapRoleFromDomain(role UserRole) authv1.UserRole {
	switch role {
	case User:
		return authv1.UserRole_USER_ROLE_USER
	case Admin:
		return authv1.UserRole_USER_ROLE_ADMIN
	default:
		return authv1.UserRole_USER_ROLE_UNSPECIFIED
	}
}
