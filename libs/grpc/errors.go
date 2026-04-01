package grpcUtils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	BodyIsRequired = status.Error(codes.InvalidArgument, "request body is required")
)
