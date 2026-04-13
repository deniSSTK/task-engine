package grpcAuth

import (
	"context"

	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
)

type AuthVerifier interface {
	Verify(ctx context.Context) (*authv1.AuthUser, error)
}
