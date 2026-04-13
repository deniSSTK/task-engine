package authGrpc

import (
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	grpcUtils "github.com/deniSSTK/task-engine/libs/grpc"
)

type Server struct {
	grpcServer *grpcUtils.GrpcServer

	authHandler *Handler
}

func NewServer(grpcServer *grpcUtils.GrpcServer, authHandler *Handler) *Server {
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)

	return &Server{grpcServer, authHandler}
}
