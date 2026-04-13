package grpcAuth

import (
	"context"

	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	"github.com/deniSSTK/task-engine/libs/logger"
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

func (rv *RemoteAuthVerifier) Verify(ctx context.Context) (*authv1.AuthUser, error) {
	resp, err := rv.authClient.Verify(ctx, &authv1.VerifyRequest{})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, MissingAuthUser
	}

	return resp.User, nil
}
