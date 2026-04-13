package app

import (
	authApp "github.com/deniSSTK/task-engine/auth-service/internal/app/auth"
	"github.com/deniSSTK/task-engine/auth-service/internal/delivery"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure"
	"go.uber.org/fx"
)

var FxApp = fx.New(
	infrastructure.Module,
	delivery.Module,

	authApp.Module,
)
