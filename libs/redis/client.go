package redis

import (
	"context"
	"libs/env"
	"libs/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(lc fx.Lifecycle, cfg *env.DefConfig, log *logger.Logger) *Redis {
	redisLog := log.Named("Redis")

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		redisLog.Fatal(FailedToConnectToRedis, zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	redisLog.Info("connected to redis",
		zap.String("addr", cfg.Redis.Addr),
		zap.Int("db", cfg.Redis.DB),
	)

	return &Redis{client}
}

func (r *Redis) Client() *redis.Client {
	return r.client
}
