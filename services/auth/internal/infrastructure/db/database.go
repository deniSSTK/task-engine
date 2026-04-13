package db

import (
	"github.com/deniSSTK/task-engine/auth-service/ent"
	"github.com/deniSSTK/task-engine/auth-service/internal/infrastructure/config"
	"github.com/deniSSTK/task-engine/libs/logger"
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
