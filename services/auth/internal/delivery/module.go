package delivery

import (
	"github.com/deniSSTK/task-engine/auth-service/internal/delivery/v1/grpc"
	"go.uber.org/fx"
)

var Module = fx.Options(
	grpcV1.Module,
)
