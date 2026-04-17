package infrastructure

import (
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/config"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/db"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/security"
	grpcErr "github.com/deniSSTK/task-engine/libs/grpc/error"
	"github.com/deniSSTK/task-engine/libs/logger"
	"github.com/deniSSTK/task-engine/libs/redis"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		config.NewConfigs,

		logger.NewLogger,
		redis.NewRedis,

		newAppErrorWrapper,
	),

	db.Module,
	security.Module,
)

func newAppErrorWrapper() *grpcErr.AppErrorWrapper {
	return grpcErr.NewAppErrorWrapper("auth.v1")
}
