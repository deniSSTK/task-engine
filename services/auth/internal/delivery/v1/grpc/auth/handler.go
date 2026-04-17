package authGrpcV1

import (
	"context"
	"errors"
	"strings"

	"buf.build/go/protovalidate"
	authService "github.com/deniSSTK/task-engine/auth-service/internal/service/auth"
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	grpcAuth "github.com/deniSSTK/task-engine/libs/auth"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	grpcErr "github.com/deniSSTK/task-engine/libs/grpc/error"
	"github.com/deniSSTK/task-engine/libs/logger"
	"github.com/deniSSTK/task-engine/libs/reasons"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type Handler struct {
	authv1.UnimplementedAuthServiceServer
	protoValidator protovalidate.Validator

	authService *authService.Service

	errorWrapper *grpcErr.AppErrorWrapper
	log          *logger.Logger
}

func NewHandler(
	protoValidator protovalidate.Validator,

	authService *authService.Service,

	errorWrapper *grpcErr.AppErrorWrapper,
	log *logger.Logger,
) *Handler {
	return &Handler{
		protoValidator: protoValidator,

		authService: authService,

		errorWrapper: errorWrapper,
		log:          log,
	}
}

// TODO: add logs

func (h *Handler) Register(
	ctx context.Context,
	dto *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	if dto == nil {
		return nil, h.errorWrapper.BodyIsRequired()
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
		return nil, h.errorWrapper.ValidationFailed(err)
	}

	resp, err := h.authService.Register(ctx, dto)
	if err != nil {
		if errors.Is(err, authService.EmailAlreadyExists) {
			return nil, h.errorWrapper.New(codes.AlreadyExists, err, reasons.EmailAlreadyExists, "email")
		}

		return nil, h.errorWrapper.New(codes.Internal, err, reasons.InternalServerError)
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
		return nil, h.errorWrapper.BodyIsRequired()
	}

	dto.Email = strings.ToLower(strings.TrimSpace(dto.Email))
	dto.Password = strings.TrimSpace(dto.Password)

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, h.errorWrapper.ValidationFailed(err)
	}

	resp, err := h.authService.Login(ctx, dto)
	if err != nil {
		if errors.Is(err, defErrors.InvalidCredentials) {
			return nil, h.errorWrapper.New(codes.Unauthenticated, err, reasons.InvalidCredentials)
		}

		return nil, h.errorWrapper.InternalServerError(err)
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
		return nil, h.errorWrapper.BodyIsRequired()
	}

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, h.errorWrapper.ValidationFailed(err)
	}

	resp, err := h.authService.Refresh(ctx, dto)
	if err != nil {
		return nil, h.errorWrapper.New(codes.Unauthenticated, err, reasons.AuthenticationFailed)
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
		return nil, h.errorWrapper.New(codes.Unauthenticated, err, reasons.AuthenticationFailed)
	}

	payload, err := h.authService.Me(ctx, token)
	if err != nil {
		log.Error(err.Error())
		return nil, h.errorWrapper.New(codes.Unauthenticated, err, reasons.AuthenticationFailed)
	}

	role, ok := userDomain.MapRoleFromDomain(payload.Role)
	if !ok {
		return nil,
			h.errorWrapper.Unauthenticated()
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
		return nil, h.errorWrapper.BodyIsRequired()
	}

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, h.errorWrapper.ValidationFailed(err)
	}

	id, err := uuid.Parse(dto.Id)
	if err != nil {
		return nil, h.errorWrapper.New(codes.InvalidArgument, err, reasons.InvalidArgument)
	}

	if err = h.authService.VerifyById(ctx, id); err != nil {
		if errors.Is(err, defErrors.NotFound) {
			return nil, h.errorWrapper.NotFound()
		}

		return nil, h.errorWrapper.InternalServerError(err)
	}

	return &authv1.VerifyByIdResponse{}, nil
}

func (h *Handler) UpdateUser(
	ctx context.Context,
	dto *authv1.UpdateUserRequest,
) (*authv1.UpdateUserResponse, error) {
	if dto == nil {
		return nil, h.errorWrapper.BodyIsRequired()
	}

	if err := h.protoValidator.Validate(dto); err != nil {
		return nil, h.errorWrapper.ValidationFailed(err)
	}

	authUser, ok := grpcAuth.AuthUserFromContext(ctx)
	if !ok {
		return nil, h.errorWrapper.Unauthenticated()
	}

	rawUser, err := h.authService.UpdateUser(ctx, authUser.Id, dto)
	if err != nil {
		return nil, h.errorWrapper.InternalServerError(err)
	}

	user, ok := MapDomainUserToProtoUser(rawUser)
	if !ok {
		return nil, h.errorWrapper.InternalServerError(defErrors.FailedMap)
	}

	return &authv1.UpdateUserResponse{
		User: user,
	}, nil
}
