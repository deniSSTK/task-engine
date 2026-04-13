package delivery

import (
	"github.com/deniSSTK/task-engine/auth-service/internal/delivery/grpc"
	"go.uber.org/fx"
)

var Module = fx.Options(
	grpc.Module,
)
