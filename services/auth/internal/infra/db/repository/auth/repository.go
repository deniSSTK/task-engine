package authRepo

import (
	"auth-service/ent"
	"auth-service/ent/user"
	"auth-service/internal/infra/db"
	txUtils "auth-service/utils/tx-utils"
	"context"
	userDomain "libs/user"
	"time"

	proto "proto/auth"

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

func (r *Repository) GetPasswordHashByEmail(
	ctx context.Context,
	email string,
) (string, error) {
	client := txUtils.FromContext(ctx, r.client)

	return client.User.
		Query().
		Where(user.EmailEQ(email)).
		Select(user.FieldPasswordHash).
		String(ctx)
}

func (r *Repository) UpdateUserLastLoginAtByEmail(
	ctx context.Context,
	email string,
) error {
	client := txUtils.FromContext(ctx, r.client)

	return client.User.
		Update().
		Where(user.EmailEQ(email)).
		SetLastLoginAt(time.Now()).
		Exec(ctx)
}

func (r *Repository) GetUserIdAndRoleByEmail(
	ctx context.Context,
	email string,
) (GetUserIdAndRoleByEmailRes, error) {
	client := txUtils.FromContext(ctx, r.client)

	var res GetUserIdAndRoleByEmailRes

	if err := client.User.
		Query().
		Where(user.EmailEQ(email)).
		Select(
			user.FieldID,
			user.FieldRole,
		).
		Scan(ctx, &res); err != nil {
		return GetUserIdAndRoleByEmailRes{}, err
	}

	return res, nil
}
