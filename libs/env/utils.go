package env

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func GetEnv(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

func EnvMust(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Env key " + key + " is missing")
	}
	return value
}

func EnvMustToStringList(key string) []string {
	value := os.Getenv(key)
	if value == "" {
		panic("Env key " + key + " is missing")
	}
	return strings.Split(value, ",")
}

func EnvMustDuration(key string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		panic("Env key " + key + " is missing")
	}

	ttl, err := time.ParseDuration(value)
	if err != nil {
		panic("failed to parse time duration" + key)
	}

	return ttl
}

func ServiceEnvPath() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic(FailedToResolveConfigPath)
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", ".env"))
}
