package grpcUtils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	BodyIsRequired = status.Error(codes.InvalidArgument, "request body is required")
)

func FieldIsRequired(field string) error {
	return status.Error(codes.InvalidArgument, field+" is required")
}
