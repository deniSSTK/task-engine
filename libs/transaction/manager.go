package transaction

import "context"

type Tx interface {
	Commit() error
	Rollback() error
}

type Manager[TTx Tx] struct {
	begin       func(ctx context.Context) (TTx, error)
	intoContext func(ctx context.Context, tx TTx) context.Context
}

func NewManager[TTx Tx](
	begin func(ctx context.Context) (TTx, error),
	intoContext func(ctx context.Context, tx TTx) context.Context,
) *Manager[TTx] {
	return &Manager[TTx]{
		begin:       begin,
		intoContext: intoContext,
	}
}

func (tm *Manager[TTx]) WithTransaction(
	ctx context.Context,
	fn func(txCtx context.Context) error,
) (err error) {
	tx, err := tm.begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}

		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txCtx := tm.intoContext(ctx, tx)

	if err = fn(txCtx); err != nil {
		return err
	}

	return tx.Commit()
}
