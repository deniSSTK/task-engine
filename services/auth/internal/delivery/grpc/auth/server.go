package authGrpc

import (
	"auth-service/internal/delivery/grpc"

	proto "proto/auth"
)

type Server struct {
	grpcServer  *grpc.GrpcServer
	authHandler *Handler
}

func NewServer(grpcServer *grpc.GrpcServer, authHandler *Handler) *Server {
	proto.RegisterAuthServiceServer(grpcServer, authHandler)

	return &Server{grpcServer, authHandler}
}
