package delivery

import (
	"auth-service/internal/delivery/grpc"

	"go.uber.org/fx"
)

var Module = fx.Options(
	grpc.Module,
)
