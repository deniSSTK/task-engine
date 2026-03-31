package authApp

import (
	authGrpc "auth-service/internal/delivery/grpc/auth"
	"auth-service/internal/infra/db/repository/auth"
	authService "auth-service/internal/service/auth"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		authRepo.NewRepository,
		authService.NewService,

		authGrpc.NewServer,
		authGrpc.NewHandler,
	),
)
