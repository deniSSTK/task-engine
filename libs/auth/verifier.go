package grpcAuth

import (
	"context"

	userDomain "github.com/deniSSTK/task-engine/libs/user"
)

type AuthVerifier interface {
	Verify(ctx context.Context) (*userDomain.AuthUser, error)
}
