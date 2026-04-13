package grpcUtils

import (
	"context"
	"net"

	"github.com/deniSSTK/task-engine/libs/env"
	"github.com/deniSSTK/task-engine/libs/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	*grpc.Server
}

type GrpcServerParams struct {
	fx.In

	Lifecycle fx.Lifecycle

	DefCfg *env.DefConfig
	Log    *logger.Logger

	UnaryInterceptors []grpc.UnaryServerInterceptor `group:"grpc_unary_interceptors"`
}

func NewGrpcServer(params GrpcServerParams) *GrpcServer {
	var opts []grpc.ServerOption

	if len(params.UnaryInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(params.UnaryInterceptors...))
	}

	grpcServer := &GrpcServer{
		Server: grpc.NewServer(),
	}

	lis, err := net.Listen("tcp", ":"+params.DefCfg.AppPort)
	if err != nil {
		panic(err)
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err = grpcServer.Serve(lis); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})

	return grpcServer
}
