package authRepo

import (
	"time"

	userDomain "github.com/deniSSTK/task-engine/libs/user"
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
