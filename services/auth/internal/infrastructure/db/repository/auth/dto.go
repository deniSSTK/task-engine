package authRepo

import (
	userDomain "libs/user"
	"time"

	"github.com/google/uuid"
)

type GetUserIdAndRoleByEmailDto struct {
	Id   uuid.UUID
	Role userDomain.UserRole
}

type GetUserStatusDto struct {
	Status    userDomain.UserStatus
	DeletedAt *time.Time
	Role      userDomain.UserRole
}
