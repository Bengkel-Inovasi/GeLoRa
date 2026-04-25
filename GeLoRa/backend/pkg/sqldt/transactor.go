package sqldt

import (
	"context"
	"database/sql"
)

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

func (t *Transactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if getTx(ctx) != nil {
		return fn(ctx)
	}

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	wrapped := &txImpl{tx: tx}
	ctx = putTx(ctx, wrapped)

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
