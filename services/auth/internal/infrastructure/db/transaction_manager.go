package db

import (
	"context"

	"github.com/deniSSTK/task-engine/auth-service/ent"
	txUtils "github.com/deniSSTK/task-engine/auth-service/utils/tx-utils"
	"github.com/deniSSTK/task-engine/libs/transaction"
)

func newTransactionManager(db *Database) *transaction.Manager[*ent.Tx] {
	return transaction.NewManager(
		func(ctx context.Context) (*ent.Tx, error) {
			return db.Client().Tx(ctx)
		},
		txUtils.WithTx,
	)
}
