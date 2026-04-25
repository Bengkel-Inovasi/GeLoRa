package sqldt

import (
	"context"
	"database/sql"
)

type Sqldt interface {
	Exec(ctx context.Context, query string, args ...any) (res sql.Result, err error)
	Query(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRow(ctx context.Context, query string, args ...any) (row *sql.Row)
	WithTx(ctx context.Context, fn func(context.Context) error) (err error)
}
