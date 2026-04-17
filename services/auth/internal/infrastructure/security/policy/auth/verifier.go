package authPolicy

import (
	"context"

	"github.com/deniSSTK/task-engine/auth-service/internal/delivery/grpc/auth"
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	grpcAuth "github.com/deniSSTK/task-engine/libs/auth"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	"github.com/deniSSTK/task-engine/libs/logger"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
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

func (lv *LocalAuthVerifier) Verify(ctx context.Context) (*userDomain.AuthUser, error) {
	res, err := lv.authHandler.Me(ctx, &authv1.MeRequest{})
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, grpcAuth.MissingAuthUser
	}

	user, ok := userDomain.MapAuthUserFromProtoAuthUser(res.User)
	if !ok {
		return nil, defErrors.UserUnauthenticated
	}

	return user, nil
}
