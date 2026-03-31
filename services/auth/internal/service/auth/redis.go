package authService

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserSessionCachePayload struct {
	UserId       uuid.UUID `redis:"user_id"`
	RefreshToken string    `redis:"refresh_token"`

	TTL time.Duration `redis:"-"`
	//IsRevoked    bool
}

func (s *Service) saveUserSessionCache(
	ctx context.Context,
	payload *UserSessionCachePayload,
) error {
	key := s.buildUserSessionCacheKey(payload.UserId.String())

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

func (s *Service) buildUserSessionCacheKey(userId string) string {
	return fmt.Sprintf("session:%s:%s", userId)
}

func (s *Service) isSessionExists(
	ctx context.Context,
	userId string,
) (bool, error) {
	key := s.buildUserSessionCacheKey(userId)

	val, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		s.log.Error(FailedToVerifySessionExistingInCache.Error(), zap.Error(err))
		return false, FailedToVerifySessionExistingInCache
	}

	return val > 0, nil
}
