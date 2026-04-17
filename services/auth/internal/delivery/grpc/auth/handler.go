package authGrpc

import (
	"context"
	"errors"
	"strings"

	"buf.build/go/protovalidate"
	authService "github.com/deniSSTK/task-engine/auth-service/internal/service/auth"
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	grpcAuth "github.com/deniSSTK/task-engine/libs/auth"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	grpcUtils "github.com/deniSSTK/task-engine/libs/grpc"
	"github.com/deniSSTK/task-engine/libs/logger"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	authv1.UnimplementedAuthServiceServer
	protoValidator protovalidate.Validator

	authService *authService.Service

	log *logger.Logger
}

func NewHandler(
	protoValidator protovalidate.Validator,

	authService *authService.Service,

	log *logger.Logger,
) *Handler {
	return &Handler{
		protoValidator: protoValidator,

		authService: authService,

		log: log,
	}
}

// TODO: add error codes

func (h *Handler) Register(
	ctx context.Context,
	dto *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
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

	return &authv1.RegisterResponse{
		Tokens: MapTokenPairToProtoTokenDetail(resp),
	}, nil
}

func (h *Handler) Login(
	ctx context.Context,
	dto *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
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

	return &authv1.LoginResponse{
		Tokens: MapTokenPairToProtoTokenDetail(resp),
	}, nil
}

func (h *Handler) Refresh(
	ctx context.Context,
	dto *authv1.RefreshRequest,
) (*authv1.RefreshResponse, error) {
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

	return &authv1.RefreshResponse{
		Tokens: MapTokenPairToProtoTokenDetail(resp),
	}, nil
}

func (h *Handler) Me(
	ctx context.Context,
	_ *authv1.MeRequest,
) (*authv1.MeResponse, error) {
	log := h.log.Named("Me")

	token, err := grpcAuth.ExtractAuthToken(ctx, h.log)
	if err != nil {
		log.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	payload, err := h.authService.Me(ctx, token)
	if err != nil {
		log.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	role, ok := userDomain.MapRoleFromDomain(payload.Role)
	if !ok {
		log.Error(defErrors.UserUnauthenticated.Error(),
			zap.String("reason", InvalidRole.Error()),
		)
		return nil, status.Error(codes.Unauthenticated, defErrors.UserUnauthenticated.Error())
	}

	return &authv1.MeResponse{
		User: &authv1.AuthUser{
			Id:   payload.UserId.String(),
			Role: role,
		},
	}, nil
}

func (h *Handler) VerifyById(
	ctx context.Context,
	dto *authv1.VerifyByIdRequest,
) (*authv1.VerifyByIdResponse, error) {
	if dto == nil {
		return nil, grpcUtils.BodyIsRequired
	}

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := uuid.Parse(dto.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = h.authService.VerifyById(ctx, id); err != nil {
		// TODO: change error code to Unauthenticated OR NotFound
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &authv1.VerifyByIdResponse{}, nil
}

func (h *Handler) UpdateUser(
	ctx context.Context,
	dto *authv1.UpdateUserRequest,
) (*authv1.UpdateUserResponse, error) {
	if dto == nil {
		return nil, grpcUtils.BodyIsRequired
	}

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	authUser, ok := grpcAuth.AuthUserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, defErrors.UserUnauthenticated.Error())
	}

	rawUser, err := h.authService.UpdateUser(ctx, authUser.Id, dto)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user, ok := MapDomainUserToProtoUser(rawUser)
	if !ok {
		return nil, status.Error(codes.Internal, defErrors.FailedMap.Error())
	}

	return &authv1.UpdateUserResponse{
		User: user,
	}, nil
}
