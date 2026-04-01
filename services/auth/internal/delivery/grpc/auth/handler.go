package authGrpc

import (
	authService "auth-service/internal/service/auth"
	"context"
	"errors"
	grpcUtils "libs/grpc"
	proto "proto/proto/auth/v1"
	"strings"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	proto.UnimplementedAuthServiceServer
	authService *authService.Service

	protoValidator protovalidate.Validator
}

func NewHandler(
	authService *authService.Service,
	protoValidator protovalidate.Validator,
) *Handler {
	return &Handler{
		authService:    authService,
		protoValidator: protoValidator,
	}
}

// TODO: add error codes

func (h *Handler) Register(ctx context.Context, dto *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if dto == nil {
		return nil, grpcUtils.BodyIsRequired
	}

	dto.Email = strings.ToLower(strings.TrimSpace(dto.Email))
	dto.Password = strings.TrimSpace(dto.Password)
	dto.Name = strings.TrimSpace(dto.Name)

	if dto.SecondName != nil {
		trimmed := strings.TrimSpace(*dto.SecondName)
		if trimmed == "" {
			dto.SecondName = nil
		} else {
			dto.SecondName = &trimmed
		}
	}

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := h.authService.Register(ctx, dto)
	if err != nil {
		if errors.Is(err, authService.EmailAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.RegisterResponse{
		Tokens: MapTokenPairToProtoTokenDetail(resp),
	}, nil
}

func (h *Handler) Login(ctx context.Context, dto *proto.LoginRequest) (*proto.LoginResponse, error) {
	if dto == nil {
		return nil, grpcUtils.BodyIsRequired
	}

	dto.Email = strings.ToLower(strings.TrimSpace(dto.Email))
	dto.Password = strings.TrimSpace(dto.Password)

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := h.authService.Login(ctx, dto)
	if err != nil {
		if errors.Is(err, authService.InvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.LoginResponse{
		Tokens: MapTokenPairToProtoTokenDetail(resp),
	}, nil
}

func (h *Handler) Refresh(ctx context.Context, dto *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	if dto == nil {
		return nil, grpcUtils.BodyIsRequired
	}

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := h.authService.Refresh(ctx, dto)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &proto.RefreshResponse{
		Tokens: MapTokenPairToProtoTokenDetail(resp),
	}, nil
}
