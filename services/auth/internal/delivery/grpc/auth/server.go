package authGrpc

import (
	grpcUtils "libs/grpc"
	proto "proto/proto/auth/v1"
)

type Server struct {
	grpcServer  *grpcUtils.GrpcServer
	authHandler *Handler
}

func NewServer(grpcServer *grpcUtils.GrpcServer, authHandler *Handler) *Server {
	proto.RegisterAuthServiceServer(grpcServer, authHandler)

	return &Server{grpcServer, authHandler}
}
