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
	DeletedAt *time.Time `sql:"deleted_at"`
	Role      userDomain.UserRole
}

type UpdateUser struct {
	Id         uuid.UUID
	Name       *string
	SecondName **string
}

type CreateUserDto struct {
	Email        string
	PasswordHash string
	Name         string
	SecondName   *string
}

type CreateUserSessionDto struct {
	RefreshToken string
	Ip           *string
	UserAgent    *string
	ExpiresAt    time.Time
	UserId       uuid.UUID
}
