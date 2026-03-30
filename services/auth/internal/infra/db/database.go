package db

import (
	"auth-service/ent"
	"auth-service/internal/infra/config"
	"libs/logger"
)

type Database struct {
	client *ent.Client
	log    *logger.Logger
}

func NewDatabase(cfg *config.Config, log *logger.Logger) *Database {
	dbLog := log.Named("Database")

	client := newClient(cfg, dbLog)

	return &Database{client: client, log: dbLog}
}
