package app

import (
	authApp "auth-service/internal/app/auth"
	"auth-service/internal/delivery"
	"auth-service/internal/infra"

	"go.uber.org/fx"
)

var FxApp = fx.New(
	infra.Module,
	delivery.Module,

	authApp.Module,
)
