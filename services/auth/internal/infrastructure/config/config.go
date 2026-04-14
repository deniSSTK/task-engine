package config

import (
	"time"

	"github.com/deniSSTK/task-engine/libs/env"
	"go.uber.org/fx"
)

type jwtConfig struct {
	Secret          string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Config struct {
	*env.ServiceConfig

	JWT jwtConfig
}

type ConfigsOut struct {
	fx.Out

	*Config
	*env.DefConfig
	*env.ServiceConfig
}

func NewConfigs() ConfigsOut {
	defCfg := env.NewDefConfig("AUTH_PORT", env.ServiceEnvPath())

	serviceConfig := env.NewServiceConfig(defCfg)

	config := &Config{
		ServiceConfig: env.NewServiceConfig(defCfg),

		JWT: buildJwtConfig(),
	}

	return ConfigsOut{
		Config:        config,
		DefConfig:     defCfg,
		ServiceConfig: serviceConfig,
	}
}

func buildJwtConfig() jwtConfig {
	return jwtConfig{
		Secret:          env.EnvMust("JWT_SECRET"),
		Issuer:          env.GetEnv("JWT_ISSUER", "auth-service"),
		AccessTokenTTL:  env.EnvMustDuration("JWT_ACCESS_TOKEN_TTL"),
		RefreshTokenTTL: env.EnvMustDuration("JWT_REFRESH_TOKEN_TTL"),
	}
}
