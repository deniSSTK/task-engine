package db

import (
	"auth-service/ent"
	txUtils "auth-service/utils/tx-utils"
	"context"
	"libs/transaction"
)

func newTransactionManager(db *Database) *transaction.Manager[*ent.Tx] {
	return transaction.NewManager(
		func(ctx context.Context) (*ent.Tx, error) {
			return db.Client().Tx(ctx)
		},
		txUtils.WithTx,
	)
}
