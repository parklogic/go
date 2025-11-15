package mysql

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"

	"github.com/parklogic/go/database"
	"github.com/parklogic/go/database/transaction"
)

func Exec(ctx context.Context, stmt Statement) (err error) {
	var db qrm.Executable

	db = transaction.Ctx(ctx)

	if db == (*sql.Tx)(nil) {
		db, err = database.Ctx(ctx)
		if err != nil {
			return err
		}
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
	var db qrm.Executable

	db = transaction.Ctx(ctx)

	if db == (*sql.Tx)(nil) {
		db, err = database.Ctx(ctx)
		if err != nil {
			return 0, err
		}
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
	var db qrm.Queryable

	db = transaction.Ctx(ctx)

	if db == (*sql.Tx)(nil) {
		db, err = database.Ctx(ctx)
		if err != nil {
			return err
		}
	}

	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		return database.TranslateError(err)
	}

	return nil
}
