package mysql

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/qrm"

	"github.com/parklogic/go/database"
	"github.com/parklogic/go/database/transaction"
	"github.com/parklogic/go/pagination"
)

func Exec(ctx context.Context, stmt Statement) (err error) {
	var db qrm.Executable

	db = transaction.Ctx(ctx)

	if db == (*sql.Tx)(nil) {
		if db, err = database.Ctx(ctx); err != nil {
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
		if db, err = database.Ctx(ctx); err != nil {
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

func PaginateStatement(ctx context.Context, stmt SelectStatement, sortableColumns ColumnList) (SelectStatement, error) {
	paging := pagination.Ctx(ctx)

	if !paging.Enabled() {
		return stmt, nil
	}

	countStmt := SELECT(COUNT(STAR).AS("count")).FROM(stmt.AsTable("c"))

	var count struct{ Count int64 }
	if err := Query(ctx, countStmt, &count); err != nil {
		return nil, err
	}

	paging.Size = count.Count

	if count.Count == 0 {
		return SELECT(NULL), nil
	}

	limit := paging.Limit()
	if limit >= 0 {
		stmt = stmt.LIMIT(limit)
	}

	offset := paging.Offset()
	if offset > 0 {
		stmt = stmt.OFFSET(offset)
	}

	sortKeys := paging.SortKeys()
	if len(sortKeys) == 0 {
		return stmt, nil
	}

	columns := make(map[string]Column, len(sortableColumns))
	for _, c := range sortableColumns {
		columns[c.Name()] = c
	}

	for _, sort := range sortKeys {
		column, sortable := columns[sort.Key]
		if !sortable {
			return nil, pagination.ErrInvalidSortKey
		}

		switch sort.Order {
		case "":
			fallthrough
		case "asc":
			stmt = stmt.ORDER_BY(column.ASC())
		case "desc":
			stmt = stmt.ORDER_BY(column.DESC())
		default:
			return nil, pagination.ErrInvalidSortOrder
		}
	}

	return stmt, nil
}

func Query(ctx context.Context, stmt Statement, dest any) (err error) {
	var db qrm.Queryable

	db = transaction.Ctx(ctx)

	if db == (*sql.Tx)(nil) {
		if db, err = database.Ctx(ctx); err != nil {
			return err
		}
	}

	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		return database.TranslateError(err)
	}

	return nil
}
