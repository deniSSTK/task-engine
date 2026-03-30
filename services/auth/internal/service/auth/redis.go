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
	UserId       uuid.UUID     `redis:"user_id"`
	RefreshToken string        `redis:"refresh_token"`
	ExpiredAt    time.Duration `redis:"expired_at"`
	//IsRevoked    bool
}

func (s *Service) saveUserSessionCache(
	ctx context.Context,
	payload *UserSessionCachePayload,
) error {
	key := s.buildUserSessionCacheKey(payload.UserId.String())

	_, err := s.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, payload)
		pipe.Expire(ctx, key, payload.ExpiredAt)
		return nil
	})

	if err != nil {
		s.log.Error(FailedToSaveUserSessionCache.Error(), zap.Error(err))
		return FailedToSaveUserSessionCache
	}

	return nil
}

func (s *Service) buildUserSessionCacheKey(userId string) string {
	return fmt.Sprintf("session:%s", userId)
}
