package authApp

import (
	authGrpc2 "github.com/deniSSTK/task-engine/auth-service/internal/delivery/v1/grpc/auth"
	authRepo "github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/db/repository/auth"
	authService "github.com/deniSSTK/task-engine/auth-service/internal/service/auth"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		authRepo.NewEntRepository,
		authService.NewService,

		//authGrpc.NewServer,
		authGrpc2.NewHandler,
	),

	fx.Invoke(authGrpc2.NewServer),
)
