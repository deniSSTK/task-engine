package authGrpc

import (
	authService "auth-service/internal/service/auth"
	"context"
	"errors"

	proto "github.com/deniSSTK/task-engine/gen/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	proto.UnimplementedAuthServiceServer
	authService *authService.Service
}

func NewHandler(authService *authService.Service) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Register(ctx context.Context, dto *proto.RegisterRequest) (*proto.TokensResponse, error) {
	if dto == nil {
		return nil, status.Error(codes.InvalidArgument, "request body is required")
	}

	if dto.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if dto.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if dto.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	resp, err := h.authService.Register(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, authService.EmailAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return resp, nil
}

func (h *Handler) Login(context.Context, *proto.LoginRequest) (*proto.TokensResponse, error) {

}

func (h *Handler) Refresh(context.Context, *proto.RefreshRequest) (*proto.TokensResponse, error) {

}
