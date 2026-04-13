package authPolicy

import (
	authv1 "github.com/deniSSTK/task-engine/gen/proto/auth/v1"
	grpcAuth "github.com/deniSSTK/task-engine/libs/auth"
	"github.com/deniSSTK/task-engine/libs/logger"
)

func NewAuthPolicyRegistry(log *logger.Logger) *grpcAuth.PolicyRegistry {
	registry := grpcAuth.NewPolicyRegistry(
		log,
		authv1.File_auth_v1_auth_proto,
	)

	return registry
}
