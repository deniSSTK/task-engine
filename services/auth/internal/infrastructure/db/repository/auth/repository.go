package authRepo

import (
	"context"
	"time"

	"github.com/deniSSTK/task-engine/auth-service/ent"
	"github.com/deniSSTK/task-engine/auth-service/ent/user"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/db"
	userMapper "github.com/deniSSTK/task-engine/auth-service/internal/mappers/user"
	txUtils "github.com/deniSSTK/task-engine/auth-service/utils/tx-utils"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
	"github.com/google/uuid"
)

type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(db *db.Database) Repository {
	return &EntRepository{client: db.Client()}
}

func (r *EntRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	client := txUtils.FromContext(ctx, r.client)

	return client.User.
		Query().
		Where(user.EmailEQ(email)).
		Exist(ctx)
}

func (r *EntRepository) CreateUser(
	ctx context.Context,
	dto *CreateUserDto,
) (uuid.UUID, userDomain.UserRole, error) {
	client := txUtils.FromContext(ctx, r.client)

	createdUser, err := client.User.
		Create().
		SetName(dto.Name).
		SetNillableSecondName(dto.SecondName).
		SetEmail(dto.Email).
		SetPasswordHash(dto.PasswordHash).
		SetLastLoginAt(time.Now()).
		Save(ctx)

	if err != nil {
		return uuid.Nil, userDomain.RoleUser, err
	}

	return createdUser.ID, userDomain.UserRole(createdUser.Role), nil
}

func (r *EntRepository) GetPasswordHashByEmail(
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

func (r *EntRepository) UpdateUserLastLoginAtByEmail(
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

func (r *EntRepository) GetUserIdAndRoleByEmail(
	ctx context.Context,
	email string,
) (GetUserIdAndRoleByEmailDto, error) {
	client := txUtils.FromContext(ctx, r.client)

	var res []GetUserIdAndRoleByEmailDto

	if err := client.User.
		Query().
		Where(user.EmailEQ(email)).
		Select(
			user.FieldID,
			user.FieldRole,
		).
		Scan(ctx, &res); err != nil {
		return GetUserIdAndRoleByEmailDto{}, err
	}

	if len(res) == 0 {
		return GetUserIdAndRoleByEmailDto{}, defErrors.NotFound
	}

	return res[0], nil
}

func (r *EntRepository) GetUserStatusDto(
	ctx context.Context,
	userId uuid.UUID,
) (GetUserStatusDto, error) {
	client := txUtils.FromContext(ctx, r.client)

	var res []GetUserStatusDto

	if err := client.User.
		Query().
		Where(user.IDEQ(userId)).
		Select(
			user.FieldStatus,
			user.FieldDeletedAt,
			user.FieldRole,
		).
		Scan(ctx, &res); err != nil {
		return GetUserStatusDto{}, err
	}

	if len(res) == 0 {
		return GetUserStatusDto{}, defErrors.NotFound
	}

	return res[0], nil
}

func (r *EntRepository) UpdateUser(
	ctx context.Context,
	dto *UpdateUser,
) (*userDomain.User, error) {
	client := txUtils.FromContext(ctx, r.client)

	update := client.User.UpdateOneID(dto.Id)

	if dto.Name != nil {
		update.SetName(*dto.Name)
	}

	if dto.SecondName != nil {
		if *dto.SecondName != nil {
			update.ClearSecondName()
		} else {
			update.SetSecondName(**dto.SecondName)
		}
	}

	rawUser, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	return userMapper.MapEntUserToDomain(rawUser), nil
}
