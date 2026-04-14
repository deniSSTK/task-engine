package env

import (
	"fmt"
	"strconv"
)

type redisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

type ServiceConfig struct {
	*DefConfig
	DBUrl string
	Redis redisConfig
}

func NewServiceConfig(defCfg *DefConfig) *ServiceConfig {
	return &ServiceConfig{
		DefConfig: defCfg,
		DBUrl:     buildDatabaseURL(),
		Redis: redisConfig{
			Addr:     buildRedisAddr(),
			Username: GetEnv("REDIS_USERNAME", ""),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       buildRedisDB(),
		},
	}
}

func buildRedisAddr() string {
	return fmt.Sprintf(
		"%s:%s",
		GetEnv("REDIS_HOST", "localhost"),
		GetEnv("REDIS_PORT", "6379"),
	)
}

func buildRedisDB() int {
	value := GetEnv("REDIS_DB", "0")
	db, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("REDIS_DB must be an integer, got %q", value))
	}

	return db
}

func buildDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		EnvMust("POSTGRES_USER"),
		EnvMust("POSTGRES_PASSWORD"),
		GetEnv("POSTGRES_HOST", "localhost"),
		GetEnv("POSTGRES_PORT", "5432"),
		EnvMust("DB_NAME"),
		GetEnv("POSTGRES_SSLMODE", "disable"),
	)
}
