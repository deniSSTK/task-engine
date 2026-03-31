package infra

import (
	"auth-service/internal/infra/config"
	"auth-service/internal/infra/db"
	"auth-service/internal/infra/security"
	"libs/logger"
	"libs/redis"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		config.NewDefConfig,

		logger.NewLogger,
		redis.NewRedis,
	),

	db.Module,
	security.Module,
)
