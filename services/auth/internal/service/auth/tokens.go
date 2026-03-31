package authService

import (
	"auth-service/internal/infra/security/jwt"
	"context"

	"go.uber.org/zap"
)

func (s *Service) generateAndStoreTokens(
	ctx context.Context,
	payload jwt.TokenPayload,
) (*jwt.TokenPair, error) {
	tokens, err := s.tokenManager.GenerateBothTokens(payload)
	if err != nil {
		s.log.Error(FailedToGenerateTokens.Error(), zap.Error(err))
		return &jwt.TokenPair{}, err
	}

	// TODO: create db session with deviceId

	cachePayload := UserSessionCachePayload{
		UserId:       payload.UserId,
		RefreshToken: tokens.RefreshToken,
		TTL:          s.config.JWT.RefreshTokenTTL,
	}

	if err = s.saveUserSessionCache(ctx, &cachePayload); err != nil {
		s.log.Error(FailedToSaveUserSessionCache.Error(), zap.Error(err))
		return &jwt.TokenPair{}, FailedToSaveUserSessionCache
	}

	return tokens, nil
}
