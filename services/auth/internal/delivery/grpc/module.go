package grpc

import (
	grpcUtils "libs/grpc"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		grpcUtils.NewGrpcServer,
	),
)
