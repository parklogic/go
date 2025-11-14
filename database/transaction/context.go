package transaction

import (
	"context"
	"database/sql"
)

type contextKey struct{}

func Ctx(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(contextKey{}).(*sql.Tx)
	if !ok {
		return nil
	}

	return tx
}

func WithContext(parent context.Context, tx *sql.Tx) (ctx context.Context) {
	return context.WithValue(parent, contextKey{}, tx)
}
