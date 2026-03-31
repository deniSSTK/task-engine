package authRepo

import (
	userDomain "libs/user"

	"github.com/google/uuid"
)

type GetUserIdAndRoleByEmailRes struct {
	Id   uuid.UUID
	Role userDomain.UserRole
}
