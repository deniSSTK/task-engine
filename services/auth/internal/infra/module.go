package infra

import (
	"auth-service/internal/infra/config"
	"auth-service/internal/infra/db"
	"libs/logger"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		logger.NewLogger,
	),

	db.Module,
)
