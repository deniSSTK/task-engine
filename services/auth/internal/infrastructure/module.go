package infrastructure

import (
	"auth-service/internal/infrastructure/config"
	"auth-service/internal/infrastructure/db"
	"auth-service/internal/infrastructure/security"
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
