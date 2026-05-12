package ports

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type TransactionWrapper interface {
	Wrap(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error
}
