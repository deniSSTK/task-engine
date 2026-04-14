package env

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type DefConfig struct {
	ENV     Env
	AppPort string
}

func NewDefConfig(appPortName string, extraEnvFiles ...string) *DefConfig {
	loadEnvFiles(extraEnvFiles...)

	return &DefConfig{
		AppPort: EnvMust(appPortName),
		ENV:     Env(EnvMust("ENV")),
	}
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
