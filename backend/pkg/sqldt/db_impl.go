package sqldt

import (
	"context"
	"database/sql"
)

type dbImpl struct {
	db *sql.DB
}

func NewSqldt(db *sql.DB) Sqldt {
	return &dbImpl{
		db: db,
	}
}

func (d *dbImpl) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.Exec(ctx, query, args...)
	}
	return d.db.ExecContext(ctx, query, args...)
}

func (d *dbImpl) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.Query(ctx, query, args...)
	}
	return d.db.QueryContext(ctx, query, args...)
}

func (d *dbImpl) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	if tx := getTx(ctx); tx != nil {
		return tx.QueryRow(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *dbImpl) WithTx(ctx context.Context, fn func(context.Context) error) error {
	if tx := getTx(ctx); tx != nil {
		return tx.WithTx(ctx, fn)
	}

	tx, err := d.db.BeginTx(ctx, nil)
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
