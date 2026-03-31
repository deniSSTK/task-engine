package security

import (
	"auth-service/internal/infra/security/jwt"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		jwt.NewTokenManager,
	),
)
