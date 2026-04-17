package authRepo

import (
	"time"

	userDomain "github.com/deniSSTK/task-engine/libs/user"
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
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

func (u *UpdateUser) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", u.Id.String())

	if u.Name != nil {
		enc.AddString("name", *u.Name)
	}

	if u.SecondName != nil {
		if *u.SecondName == nil {
			enc.AddString("second_name", "nil")
		} else {
			enc.AddString("second_name", **u.SecondName)
		}
	}

	return nil
}
