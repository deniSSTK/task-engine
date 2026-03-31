package main

import (
	"libs/env"
)

type Config struct {
	AppPort  string
	AuthPort string
}

func NewConfig() *Config {
	return &Config{
		AppPort:  env.EnvMust("GATEWAY_PORT"),
		AuthPort: env.EnvMust("AUTH_PORT"),
	}
}
