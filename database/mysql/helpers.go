package mysql

import (
	"context"

	. "github.com/go-jet/jet/v2/mysql"
	"github.com/parklogic/go/database"
	"github.com/parklogic/go/database/internal/dbctx"
)

func Exec(ctx context.Context, stmt Statement) (err error) {
	db, err := dbctx.DB(ctx)
	if err != nil {
		return err //nolint:wrapcheck // missing database keeps the context error's text
	}

	res, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return database.TranslateError(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}

func Insert(ctx context.Context, stmt Statement) (id int64, err error) {
	db, err := dbctx.DB(ctx)
	if err != nil {
		return 0, err //nolint:wrapcheck // missing database keeps the context error's text
	}

	res, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return 0, database.TranslateError(err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Query(ctx context.Context, stmt Statement, dest any) (err error) {
	db, err := dbctx.DB(ctx)
	if err != nil {
		return err //nolint:wrapcheck // missing database keeps the context error's text
	}

	err = stmt.QueryContext(ctx, db, dest)
	if err != nil {
		return database.TranslateError(err)
	}

	return nil
}
