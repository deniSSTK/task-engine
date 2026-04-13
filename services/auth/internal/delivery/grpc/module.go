package grpc

import (
	"buf.build/go/protovalidate"
	grpcUtils "github.com/deniSSTK/task-engine/libs/grpc"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		grpcUtils.NewGrpcServer,

		newValidator,

		grpcUtils.NewUnaryInterceptor,
	),
)

func newValidator() protovalidate.Validator {
	validator, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	return validator
}
