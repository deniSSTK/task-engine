package authPolicy

import (
	"context"

	"github.com/deniSSTK/task-engine/auth-service/internal/delivery/grpc/auth"
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	grpcAuth "github.com/deniSSTK/task-engine/libs/auth"
	"github.com/deniSSTK/task-engine/libs/logger"
)

type LocalAuthVerifier struct {
	authHandler *authGrpc.Handler

	log *logger.Logger
}

func NewLocalVerifier(
	authHandler *authGrpc.Handler,
	log *logger.Logger,
) grpcAuth.AuthVerifier {
	return &LocalAuthVerifier{
		authHandler,

		log,
	}
}

func (lv *LocalAuthVerifier) Verify(ctx context.Context) (*authv1.AuthUser, error) {
	res, err := lv.authHandler.Me(ctx, &authv1.MeRequest{})
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, grpcAuth.MissingAuthUser
	}

	return res.User, nil
}
