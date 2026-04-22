package transaction

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxRunner struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *TxRunner {
	return &TxRunner{
		pool: pool,
	}
}

func (t *TxRunner) Wrap(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback(ctx) }()

	err = fn(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
