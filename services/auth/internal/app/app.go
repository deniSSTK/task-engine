package app

import (
	authApp "auth-service/internal/app/auth"
	"auth-service/internal/delivery"
	"auth-service/internal/infrastructure"

	"go.uber.org/fx"
)

var FxApp = fx.New(
	infrastructure.Module,
	delivery.Module,

	authApp.Module,
)
