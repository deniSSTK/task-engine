package security

import (
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/security/jwt"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/security/policy/auth"
	"go.uber.org/fx"
)

var Module = fx.Options(
	authPolicy.Module,

	fx.Provide(
		jwt.NewTokenManager,
	),
)
