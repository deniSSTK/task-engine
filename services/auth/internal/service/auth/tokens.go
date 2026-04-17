package authService

import (
	"context"

	authRepo "github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/db/repository/auth"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/security/jwt"
	"go.uber.org/zap"
)

func (s *Service) generateAndStoreTokens(
	ctx context.Context,
	payload jwt.TokenPayload,
	payloadDb *authRepo.CreateUserSessionDto,
) (*jwt.TokenPair, error) {
	tokens, err := s.tokenManager.GenerateBothTokens(payload)
	if err != nil {
		s.log.Error(FailedToGenerateTokens.Error(), zap.Error(err))
		return &jwt.TokenPair{}, err
	}

	payloadDb.RefreshToken = tokens.RefreshToken // TODO: save hash instead
	payloadDb.ExpiresAt = tokens.RefreshExpiredAt

	if err = s.authRepo.CreateUserSession(ctx, payloadDb); err != nil {
		s.log.Error(FailedToSaveUserSessionInDb.Error(), zap.Error(err))
		return &jwt.TokenPair{}, FailedToSaveUserSessionInDb
	}

	cachePayload := UserSessionCachePayload{
		UserId:           payload.UserId,
		RefreshTokenHash: tokens.RefreshToken,
		AccessTokenHash:  tokens.AccessToken,
		TTL:              s.config.JWT.RefreshTokenTTL,
	}

	if err = s.saveUserSessionCache(ctx, &cachePayload); err != nil {
		s.log.Error(FailedToSaveUserSessionCache.Error(), zap.Error(err))
		return &jwt.TokenPair{}, FailedToSaveUserSessionCache
	}

	return tokens, nil
}

// TODO: validate token in cache
