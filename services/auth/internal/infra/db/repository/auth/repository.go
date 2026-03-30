package authRepo

import (
	"auth-service/ent"
	"auth-service/ent/user"
	"auth-service/internal/infra/db"
	txUtils "auth-service/utils/tx-utils"
	"context"
	userDomain "libs/user"
	"time"

	proto "github.com/deniSSTK/task-engine/gen/auth"
	"github.com/google/uuid"
)

type Repository struct {
	client *ent.Client
}

func NewRepository(db *db.Database) *Repository {
	return &Repository{client: db.Client()}
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	client := txUtils.FromContext(ctx, r.client)

	return client.User.
		Query().
		Where(user.EmailEQ(email)).
		Exist(ctx)
}

func (r *Repository) CreateUser(
	ctx context.Context,
	dto *proto.RegisterRequest,
	passwordHash string,
) (uuid.UUID, userDomain.UserRole, error) {
	client := txUtils.FromContext(ctx, r.client)

	var fullName string

	if dto.SecondName != nil {
		fullName = dto.Name + " " + *dto.SecondName
	} else {
		fullName = dto.Name
	}

	createdUser, err := client.User.
		Create().
		SetName(dto.Name).
		SetNillableSecondName(dto.SecondName).
		SetEmail(dto.Email).
		SetPasswordHash(passwordHash).
		SetFullName(fullName).
		SetLastLoginAt(time.Now()).
		Save(ctx)

	if err != nil {
		return uuid.Nil, userDomain.User, err
	}

	return createdUser.ID, userDomain.UserRole(createdUser.Role), nil
}
