package grpcAuth

import (
	"context"

	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	defErrors "github.com/deniSSTK/task-engine/libs/errors"
	"github.com/deniSSTK/task-engine/libs/logger"
	userDomain "github.com/deniSSTK/task-engine/libs/user"
)

type RemoteAuthVerifier struct {
	authClient authv1.AuthServiceClient

	log *logger.Logger
}

func NewRemoteVerifier(
	authClient authv1.AuthServiceClient,

	log *logger.Logger,
) AuthVerifier {
	return &RemoteAuthVerifier{
		authClient,

		log,
	}
}

func (rv *RemoteAuthVerifier) Verify(ctx context.Context) (*userDomain.AuthUser, error) {
	resp, err := rv.authClient.Me(ctx, &authv1.MeRequest{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, MissingAuthUser
	}

	user, ok := userDomain.MapAuthUserFromProtoAuthUser(resp.User)
	if !ok {
		return nil, defErrors.UserUnauthenticated
	}

	return user, nil
}
