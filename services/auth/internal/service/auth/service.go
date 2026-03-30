package authService

import (
	"auth-service/ent"
	"auth-service/internal/infra/config"
	"auth-service/internal/infra/db/repository/auth"
	"auth-service/internal/infra/security/jwt"
	"auth-service/utils"
	"context"
	"libs/logger"
	"libs/redis"
	"libs/transaction"

	redisClient "github.com/redis/go-redis/v9"

	proto "github.com/deniSSTK/task-engine/gen/auth"
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

func (s *Service) Register(ctx context.Context, dto *proto.RegisterRequest) (*proto.TokensResponse, error) {
	s.log.Info("Registering user")

	emailExists, err := s.authRepo.EmailExists(ctx, dto.Email)

	if err != nil {
		s.log.Error(
			FailedToValidateCredentials.Error(),
			zap.Error(err),
			zap.String("email", dto.Email),
		)
		return nil, err
	}

	if emailExists {
		s.log.Error(EmailAlreadyExists.Error())
		return nil, EmailAlreadyExists
	}

	passwordHash, err := utils.HashPassword(dto.Password)
	if err != nil {
		s.log.Error(FailedToCreateUser.Error(), zap.Error(err))
		return nil, FailedToCreateUser
	}

	var tokens *jwt.TokenPair

	if err = s.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {

		userId, role, txErr := s.authRepo.CreateUser(txCtx, dto, passwordHash)
		if txErr != nil {
			s.log.Error(FailedToCreateUser.Error(), zap.Error(txErr))
			return txErr
		}

		tokenPayload := jwt.TokenPayload{
			UserId: userId,
			Role:   role,
		}

		tokens, txErr = s.tokenManager.GenerateBothTokens(tokenPayload)
		if txErr != nil {
			s.log.Error(FailedToCreateUser.Error(), zap.Error(txErr))
			return txErr
		}

		// TODO: create db session with deviceId

		cachePayload := UserSessionCachePayload{
			UserId:       userId,
			RefreshToken: tokens.RefreshToken,
			ExpiredAt:    s.config.JWT.RefreshTokenTTL,
		}

		if txErr = s.saveUserSessionCache(txCtx, &cachePayload); txErr != nil {
			s.log.Error(FailedToCreateUser.Error(), zap.Error(err))
			return FailedToCreateUser
		}

		return nil
	}); err != nil {
		s.log.Error(FailedToCreateUser.Error(), zap.Error(err))
		return nil, FailedToCreateUser
	}

	return &proto.TokensResponse{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		RefreshExpiresAt: tokens.RefreshExpiredAt.Unix(),
	}, nil
}

func (s *Service) Login(context.Context, *proto.LoginRequest) (*proto.TokensResponse, error) {

}

func (s *Service) Refresh(context.Context, *proto.RefreshRequest) (*proto.TokensResponse, error) {

}
