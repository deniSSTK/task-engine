package authService

import (
	"context"

	"github.com/deniSSTK/task-engine/auth-service/ent"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/config"
	authRepo "github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/db/repository/auth"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/security/jwt"
	"github.com/deniSSTK/task-engine/auth-service/utils"
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	"github.com/deniSSTK/task-engine/libs/logger"
	"github.com/deniSSTK/task-engine/libs/redis"
	"github.com/deniSSTK/task-engine/libs/transaction"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
	redisClient "github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

type Service struct {
	authRepo    *authRepo.Repository
	redisClient *redisClient.Client

	tokenManager       *jwt.TokenManager
	transactionManager *transaction.Manager[*ent.Tx]

	log    *logger.Logger
	config *config.Config
}

func NewService(
	authRepo *authRepo.Repository,
	redis *redis.Redis,

	tokenManager *jwt.TokenManager,
	transactionManager *transaction.Manager[*ent.Tx],

	log *logger.Logger,
	cfg *config.Config,
) *Service {
	authServiceLog := log.Named("AuthService")

	return &Service{
		authRepo:    authRepo,
		redisClient: redis.Client(),

		tokenManager:       tokenManager,
		transactionManager: transactionManager,

		log:    authServiceLog,
		config: cfg,
	}
}

func (s *Service) Register(ctx context.Context, dto *authv1.RegisterRequest) (*jwt.TokenPair, error) {
	log := s.log.Named("Register")

	emailExists, err := s.authRepo.EmailExists(ctx, dto.Email)

	if err != nil {
		log.Error(
			FailedToValidateCredentials.Error(),
			zap.Error(err),
			zap.String("email", dto.Email),
		)
		return nil, err
	}

	if emailExists {
		log.Error(EmailAlreadyExists.Error())
		return nil, EmailAlreadyExists
	}

	passwordHash, err := utils.HashPassword(dto.Password)
	if err != nil {
		log.Error(FailedToCreateUser.Error(), zap.Error(err))
		return nil, FailedToCreateUser
	}

	var tokens *jwt.TokenPair

	if err = s.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {

		userId, role, txErr := s.authRepo.CreateUser(txCtx, dto, passwordHash)
		if txErr != nil {
			log.Error(FailedToCreateUser.Error(), zap.Error(txErr))
			return txErr
		}

		payload := jwt.TokenPayload{
			UserId: userId,
			Role:   role,
		}

		tokens, txErr = s.generateAndStoreTokens(txCtx, payload)
		if txErr != nil {
			log.Error(FailedToGenerateTokens.Error(), zap.Error(txErr))
			return FailedToGenerateTokens
		}

		return nil
	}); err != nil {
		log.Error(FailedToCreateUser.Error(), zap.Error(err))
		return nil, FailedToCreateUser
	}

	return tokens, nil
}

func (s *Service) Login(ctx context.Context, dto *authv1.LoginRequest) (*jwt.TokenPair, error) {
	log := s.log.Named("Login")
	email := dto.Email

	passwordHash, err := s.authRepo.GetPasswordHashByEmail(ctx, email)
	if err != nil {
		log.Error(FailedToValidateCredentials.Error(), zap.Error(err))
		return nil, FailedToLogin
	}

	if err = utils.CheckPassword(passwordHash, dto.Password); err != nil {
		log.Error(InvalidCredentials.Error(), zap.Error(err))
		return nil, InvalidCredentials
	}

	var tokens *jwt.TokenPair

	if err = s.transactionManager.WithTransaction(ctx, func(txCtx context.Context) (txErr error) {
		if txErr = s.authRepo.UpdateUserLastLoginAtByEmail(txCtx, email); txErr != nil {
			log.Error(FailedToUpdateUserInfo.Error(), zap.Error(txErr))
			return FailedToUpdateUserInfo
		}

		targetUser, txErr := s.authRepo.GetUserIdAndRoleByEmail(txCtx, email)
		if err != nil {
			log.Error(defErrors.FailedToGetData.Error(), zap.Error(txErr))
			return defErrors.FailedToGetData
		}

		payload := jwt.TokenPayload{
			UserId: targetUser.Id,
			Role:   targetUser.Role,
		}

		tokens, txErr = s.generateAndStoreTokens(txCtx, payload)
		if txErr != nil {
			s.log.Error(FailedToGenerateTokens.Error(), zap.Error(txErr))
			return FailedToGenerateTokens
		}

		return nil
	}); err != nil {
		return &jwt.TokenPair{}, FailedToLogin
	}

	return tokens, nil
}

func (s *Service) Refresh(ctx context.Context, dto *authv1.RefreshRequest) (*jwt.TokenPair, error) {
	log := s.log.Named("Refresh")

	tokenPayload, refreshExpiresAt, err := s.tokenManager.ParseTokenPayload(dto.RefreshToken, jwt.Refresh)
	if err != nil {
		log.Error(FailedToGetTokenPayload.Error(), zap.Error(err))
		return nil, err
	}

	sessionExists, err := s.isSessionExists(ctx, tokenPayload.UserId.String())
	if err != nil {
		log.Error(FailedToVerifySessionExistingInCache.Error(), zap.Error(err))
		return nil, err
	}

	if !sessionExists {
		log.Error(FailedToVerifySessionExistingInCache.Error())
		return nil, FailedToVerifySessionExistingInCache
	}

	accessToken, err := s.tokenManager.GenerateToken(*tokenPayload, jwt.Access)
	if err != nil {
		log.Error(
			FailedToGenerateTokens.Error(),
			zap.Error(err),
			zap.String("tokenType", string(jwt.Access)),
		)
		return nil, err
	}

	return &jwt.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     dto.RefreshToken,
		RefreshExpiredAt: refreshExpiresAt.Time,
	}, nil
}

func (s *Service) Verify(ctx context.Context, token string) (*jwt.TokenPayload, error) {
	log := s.log.Named("Authorize")

	payload, _, err := s.tokenManager.ParseTokenPayload(token, jwt.Access)
	if err != nil {
		log.Error(FailedToAuthenticateUser.Error(), zap.Error(err))
		return nil, err
	}

	// TODO: get status from cache
	status, err := s.authRepo.GetUserStatusDto(ctx, payload.UserId)
	if err != nil {
		log.Error(FailedToAuthenticateUser.Error(), zap.Error(err))
		return nil, err
	}

	switch status.Status {
	case userDomain.Active:
		break
	case userDomain.Blocked:
		log.Error(UserBlocked.Error())
		return &jwt.TokenPayload{}, UserBlocked
	default:
		log.Error(
			FailedToAuthenticateUser.Error(),
			zap.Error(UndefinedUserStatus),
			zap.String("status", string(status.Status)))
	}

	if status.DeletedAt != nil {
		log.Error(UserDeleted.Error())
		return &jwt.TokenPayload{}, UserDeleted
	}

	return &jwt.TokenPayload{
		UserId: payload.UserId,
		Role:   status.Role,
	}, nil
}
