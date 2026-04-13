package authPolicy

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewAuthPolicyRegistry,
		NewLocalVerifier,
	),
)
