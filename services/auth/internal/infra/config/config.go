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

func NewConfig(defConfig *env.DefConfig) *Config {
	return &Config{
		DefConfig: defConfig,

		JWT: buildJwtConfig(),
	}
}

func NewDefConfig() *env.DefConfig {
	return env.NewDefConfig("AUTH_PORT", serviceEnvPath())
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
