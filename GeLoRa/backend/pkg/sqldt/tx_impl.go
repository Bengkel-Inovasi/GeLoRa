package sqldt

import (
	"context"
	"database/sql"
)

type txImpl struct {
	tx *sql.Tx
}

func (t *txImpl) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *txImpl) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *txImpl) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *txImpl) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
