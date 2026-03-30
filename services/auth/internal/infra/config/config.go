package config

import (
	"path/filepath"
	"runtime"
	"time"

	"libs/env"
)

type jwtConfig struct {
	Secret          string
	Issuer          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Config struct {
	*env.DefConfig

	JWT jwtConfig
}

func NewConfig() *Config {
	return &Config{
		DefConfig: env.NewDefConfig(serviceEnvPath()),

		JWT: buildJwtConfig(),
	}
}

func serviceEnvPath() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic(env.FailedToResolveConfigPath)
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", ".env"))
}

func buildJwtConfig() jwtConfig {
	return jwtConfig{
		Secret:          env.EnvMust("JWT_SECRET"),
		Issuer:          env.GetEnv("JWT_ISSUER", "auth-service"),
		AccessTokenTTL:  env.EnvMustDuration("JWT_ACCESS_TOKEN_TTL"),
		RefreshTokenTTL: env.EnvMustDuration("JWT_REFRESH_TOKEN_TTL"),
	}
}
