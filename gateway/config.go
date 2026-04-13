package main

import (
	"github.com/deniSSTK/task-engine/libs/env"
)

type Config struct {
	AppPort  string
	AuthPort string
	AuthHost string
}

func NewConfig() *Config {
	return &Config{
		AppPort: env.EnvMust("GATEWAY_PORT"),

		AuthPort: env.EnvMust("AUTH_PORT"),
		AuthHost: env.GetEnv("AUTH_HOST", "localhost"),
	}
}
