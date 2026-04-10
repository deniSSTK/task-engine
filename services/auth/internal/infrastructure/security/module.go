package security

import (
	"auth-service/internal/infrastructure/security/jwt"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		jwt.NewTokenManager,
	),
)
