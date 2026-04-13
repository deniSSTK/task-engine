package grpcUtils

import (
	grpcAuth "github.com/deniSSTK/task-engine/libs/auth"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type UnaryInterceptorOut struct {
	fx.Out

	Interceptor grpc.UnaryServerInterceptor `group:"grpc_unary_interceptor"`
}

// TODO: rename to NewAuthInterceptor and return UnaryAuthInterceptorOut

func NewUnaryInterceptor(
	registry *grpcAuth.PolicyRegistry,
	verifier grpcAuth.AuthVerifier,
) UnaryInterceptorOut {
	return UnaryInterceptorOut{
		Interceptor: grpcAuth.UnaryAuthInterceptor(registry, verifier),
	}
}
