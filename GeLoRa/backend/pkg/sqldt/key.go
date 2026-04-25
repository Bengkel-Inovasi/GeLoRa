package sqldt

import "context"

type txKey struct{}

func getTx(ctx context.Context) Sqldt {
	if pgxdx, ok := ctx.Value(txKey{}).(Sqldt); ok {
		return pgxdx
	}
	return nil
}

func putTx(ctx context.Context, inst Sqldt) context.Context {
	return context.WithValue(ctx, txKey{}, inst)
}
