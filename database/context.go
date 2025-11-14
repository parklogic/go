package database

import (
	"context"
	"database/sql"
	"errors"
)

type contextKey struct{}

func Ctx(ctx context.Context) (*sql.DB, error) {
	db, ok := ctx.Value(contextKey{}).(*sql.DB)
	if !ok {
		return nil, errors.New("database client not found in context")
	}

	return db, nil
}

func WithContext(parent context.Context, db *sql.DB) (ctx context.Context) {
	return context.WithValue(parent, contextKey{}, db)
}
