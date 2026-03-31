package grpcUtils

import (
	"context"
	"libs/env"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	*grpc.Server
}

func NewGrpcServer(lc fx.Lifecycle, defCfg *env.DefConfig) *GrpcServer {
	grpcServer := &GrpcServer{
		Server: grpc.NewServer(),
	}

	lis, err := net.Listen("tcp", ":"+defCfg.AppPort)
	if err != nil {
		panic(err)
	}

	lc.Append(fx.Hook{
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
