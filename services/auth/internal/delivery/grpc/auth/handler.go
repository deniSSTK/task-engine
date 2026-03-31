package authGrpc

import (
	authService "auth-service/internal/service/auth"
	"context"
	"errors"
	grpcUtils "libs/grpc"

	proto "proto/auth"

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
		return nil, grpcUtils.BodyIsRequired
	}

	if dto.Email == "" {
		return nil, grpcUtils.FieldIsRequired("email")
	}

	if dto.Password == "" {
		return nil, grpcUtils.FieldIsRequired("password")
	}

	if dto.Name == "" {
		return nil, grpcUtils.FieldIsRequired("name")
	}

	resp, err := h.authService.Register(ctx, dto)
	if err != nil {
		if errors.Is(err, authService.EmailAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return MapTokenPairToProtoTokenResponse(resp), nil
}

func (h *Handler) Login(ctx context.Context, dto *proto.LoginRequest) (*proto.TokensResponse, error) {
	if dto == nil {
		return nil, grpcUtils.BodyIsRequired
	}

	if dto.Email == "" {
		return nil, grpcUtils.FieldIsRequired("email")
	}

	if dto.Password == "" {
		return nil, grpcUtils.FieldIsRequired("password")
	}

	resp, err := h.authService.Login(ctx, dto)
	if err != nil {
		if errors.Is(err, authService.InvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return MapTokenPairToProtoTokenResponse(resp), nil
}

func (h *Handler) Refresh(_ context.Context, dto *proto.RefreshRequest) (*proto.TokensResponse, error) {
	if dto == nil {
		return nil, grpcUtils.BodyIsRequired
	}

	if dto.RefreshToken == "" {
		return nil, grpcUtils.FieldIsRequired("refresh_token")
	}

	resp, err := h.authService.Refresh(dto)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return MapTokenPairToProtoTokenResponse(resp), nil
}
