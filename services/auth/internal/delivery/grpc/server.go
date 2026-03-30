package grpc

import "google.golang.org/grpc"

type GrpcServer struct {
	*grpc.Server
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		grpc.NewServer(),
	}
}
