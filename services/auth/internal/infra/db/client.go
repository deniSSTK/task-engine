package db

import (
	"auth-service/ent"
	"auth-service/internal/infra/config"
	defErrors "libs/errors"
	"libs/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func (d *Database) Client() *ent.Client {
	return d.client
}

func newClient(cfg *config.Config, log *logger.Logger) *ent.Client {
	client, err := ent.Open("postgres", cfg.DBUrl)

	if err != nil {
		log.Fatal(defErrors.FailedToCreateEntClient, zap.Error(err))
	}

	return client
}
