package authService

import (
	"context"
	"fmt"
	"time"

	"github.com/deniSSTK/task-engine/auth-service/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserSessionCachePayload struct {
	UserId           uuid.UUID `redis:"user_id"`
	AccessTokenHash  string    `redis:"access_token"`
	RefreshTokenHash string    `redis:"refresh_token_hash"`

	TTL time.Duration `redis:"-"`
	//IsRevoked    bool
}

func (s *Service) saveUserSessionCache(
	ctx context.Context,
	payload *UserSessionCachePayload,
) error {
	payload.RefreshTokenHash = utils.HashHex(payload.RefreshTokenHash)
	payload.AccessTokenHash = utils.HashHex(payload.AccessTokenHash)

	key := s.buildUserSessionCacheKey(payload.UserId)

	_, err := s.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, payload)
		pipe.Expire(ctx, key, payload.TTL)
		return nil
	})

	if err != nil {
		s.log.Error(FailedToSaveUserSessionCache.Error(), zap.Error(err))
		return FailedToSaveUserSessionCache
	}

	return nil
}

func (s *Service) buildUserSessionCacheKey(userId uuid.UUID) string {
	return fmt.Sprintf("session:%s", userId.String())
}

func (s *Service) isSessionValid(
	ctx context.Context,
	userId uuid.UUID,
	refreshToken string,
) (bool, error) {
	key := s.buildUserSessionCacheKey(userId)

	existsInDb, err := s.authRepo.IsExistsSession(ctx, userId, refreshToken)
	if err != nil {
		return false, err
	}

	if !existsInDb {
		return false, nil
	}

	var cachedSession UserSessionCachePayload

	if err := s.redisClient.HGetAll(ctx, key).Scan(&cachedSession); err != nil {
		s.log.Error(FailedToVerifySessionExistingInCache.Error(), zap.Error(err))
		return false, FailedToVerifySessionExistingInCache
	}

	return cachedSession.RefreshTokenHash == utils.HashHex(refreshToken), nil
}

func (s *Service) deleteSessionFromCache(
	ctx context.Context,
	userId uuid.UUID,
) error {
	return s.redisClient.Del(ctx, s.buildUserSessionCacheKey(userId)).Err()
}
