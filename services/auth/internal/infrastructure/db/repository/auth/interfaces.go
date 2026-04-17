package authRepo

import (
	"context"

	userDomain "github.com/deniSSTK/task-engine/libs/user"
	"github.com/google/uuid"
)

type Repository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	GetUserStatusDto(ctx context.Context, userId uuid.UUID) (GetUserStatusDto, error)
	GetPasswordHashByEmail(ctx context.Context, email string) (string, error)
	GetUserIdAndRoleByEmail(ctx context.Context, email string) (GetUserIdAndRoleByEmailDto, error)

	CreateUser(ctx context.Context, dto *CreateUserDto) (uuid.UUID, userDomain.UserRole, error)
	CreateUserSession(ctx context.Context, dto *CreateUserSessionDto) error

	UpdateUser(ctx context.Context, dto *UpdateUser) (*userDomain.User, error)
	UpdateUserLastLoginAtByEmail(ctx context.Context, email string) error

	DeleteUserSession(ctx context.Context, userId uuid.UUID) error
}
