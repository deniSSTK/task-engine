package infrastructure

import (
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/config"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/db"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/security"
	"github.com/deniSSTK/task-engine/libs/logger"
	"github.com/deniSSTK/task-engine/libs/redis"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		config.NewConfigs,

		logger.NewLogger,
		redis.NewRedis,
	),

	db.Module,
	security.Module,
)
