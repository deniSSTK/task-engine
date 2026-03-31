package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

type redisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

type DefConfig struct {
	DBUrl   string
	Redis   redisConfig
	ENV     Env
	AppPort string
}

func NewDefConfig(extraEnvFiles ...string) *DefConfig {
	loadEnvFiles(extraEnvFiles...)

	return &DefConfig{
		DBUrl:   buildDatabaseURL(),
		AppPort: EnvMust("APP_PORT"),
		Redis: redisConfig{
			Addr:     buildRedisAddr(),
			Username: GetEnv("REDIS_USERNAME", ""),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       buildRedisDB(),
		},
		ENV: Env(EnvMust("ENV")),
	}
}

func buildDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		EnvMust("POSTGRES_USER"),
		EnvMust("POSTGRES_PASSWORD"),
		EnvMust("POSTGRES_HOST"),
		EnvMust("POSTGRES_PORT"),
		EnvMust("DB_NAME"),
		EnvMust("POSTGRES_SSLMODE"),
	)
}

func buildRedisAddr() string {
	return fmt.Sprintf("%s:%s", EnvMust("REDIS_HOST"), EnvMust("REDIS_PORT"))
}

func buildRedisDB() int {
	value := GetEnv("REDIS_DB", "0")
	db, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("REDIS_DB must be an integer, got %q", value))
	}

	return db
}

func loadEnvFiles(extraEnvFiles ...string) {
	envFiles := append([]string{projectEnvPath()}, extraEnvFiles...)
	merged := make(map[string]string)

	for _, envFile := range envFiles {
		if envFile == "" {
			continue
		}

		fileEnv, err := godotenv.Read(envFile)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			panic(err)
		}

		for key, value := range fileEnv {
			merged[key] = value
		}
	}

	for key, value := range merged {
		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}

func projectEnvPath() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to resolve env package path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".env"))
}
