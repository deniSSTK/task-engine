package grpcErr

import (
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	"github.com/deniSSTK/task-engine/libs/reasons"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type AppErrorWrapper struct {
	domain string
}

type AppError struct {
	Code    codes.Code
	Message string
	Reason  reasons.Reason
	Field   string
}

func NewAppErrorWrapper(domain string) *AppErrorWrapper {
	return &AppErrorWrapper{domain}
}

func (w *AppErrorWrapper) New(
	code codes.Code,
	msg error,
	reason reasons.Reason,
	field ...string,
) error {
	appErr := &AppError{
		Code:    code,
		Message: msg.Error(),
		Reason:  reason,
	}

	if len(field) > 0 {
		appErr.Field = field[0]
	}

	return w.ToGRPC(appErr)
}

func (w *AppErrorWrapper) ToGRPC(err *AppError) error {
	st := status.New(err.Code, err.Message)

	errorInfo := &errdetails.ErrorInfo{
		Domain: w.domain,
		Reason: string(err.Reason),
	}

	var details []protoadapt.MessageV1

	details = append(details, protoadapt.MessageV1Of(errorInfo))

	if err.Field != "" {
		br := &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       err.Field,
					Description: err.Message,
				},
			},
		}
		details = append(details, protoadapt.MessageV1Of(br))
	}

	st, _ = st.WithDetails(details...)
	return st.Err()
}

func (w *AppErrorWrapper) NotFound() error {
	return w.New(codes.NotFound, defErrors.NotFound, reasons.NotFound)
}

func (w *AppErrorWrapper) BodyIsRequired() error {
	return w.New(codes.InvalidArgument, defErrors.BodyIsRequired, reasons.BodyIsRequired)
}

func (w *AppErrorWrapper) Unauthenticated(rawErr error) error {
	if rawErr != nil {
		return w.New(codes.Unauthenticated, rawErr, reasons.AuthenticationFailed)
	}

	return w.New(codes.Unauthenticated, defErrors.UserUnauthenticated, reasons.AuthenticationFailed)
}

func (w *AppErrorWrapper) PermissionDenied() error {
	return w.New(codes.PermissionDenied, defErrors.PermissionDenied, reasons.PermissionDenied)
}

func (w *AppErrorWrapper) ValidationFailed(err error) error {
	return w.New(codes.InvalidArgument, err, reasons.FailedValidation)
}

func (w *AppErrorWrapper) InternalServerError(err error) error {
	return w.New(codes.Internal, err, reasons.InternalServerError)
}
