package authApp

import (
	authGrpc "github.com/deniSSTK/task-engine/auth-service/internal/delivery/grpc/auth"
	authRepo "github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/db/repository/auth"
	authService "github.com/deniSSTK/task-engine/auth-service/internal/service/auth"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		authRepo.NewRepository,
		authService.NewService,

		//authGrpc.NewServer,
		authGrpc.NewHandler,
	),

	fx.Invoke(authGrpc.NewServer),
)
