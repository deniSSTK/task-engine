package grpc

import (
	grpcUtils "libs/grpc"

	"buf.build/go/protovalidate"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		grpcUtils.NewGrpcServer,

		newValidator,
	),
)

func newValidator() protovalidate.Validator {
	validator, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	return validator
}
