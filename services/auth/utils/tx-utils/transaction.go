package txUtils

import (
	"auth-service/ent"
	"context"
	"libs/transaction"
)

func WithTx(ctx context.Context, tx *ent.Tx) context.Context {
	return context.WithValue(ctx, transaction.TxKey{}, tx)
}

func FromContext(ctx context.Context, def *ent.Client) *ent.Client {
	tx, ok := ctx.Value(transaction.TxKey{}).(*ent.Tx)
	if ok {
		return tx.Client()
	}

	return def
}
