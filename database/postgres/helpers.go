package postgres

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/parklogic/go/database"
	"github.com/parklogic/go/database/internal/dbctx"
	"github.com/parklogic/go/pagination"
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

// Insert executes stmt. It does not return an id: the default pgx driver rejects LastInsertId.
// Callers that need the new id use Query with RETURNING.
func Insert(ctx context.Context, stmt Statement) error {
	db, err := dbctx.DB(ctx)
	if err != nil {
		return err //nolint:wrapcheck // missing database keeps the context error's text
	}

	_, err = stmt.ExecContext(ctx, db)
	if err != nil {
		return database.TranslateError(err) //nolint:wrapcheck // insert returns the translated driver error
	}

	return nil
}

// Count runs the count query for an enabled page and stores the total on Pagination.Size.
// A disabled page returns without querying. Call Count before PaginateStatement.
func Count(ctx context.Context, stmt SelectStatement) error {
	paging := pagination.Ctx(ctx)
	if !paging.Enabled() {
		return nil
	}

	countStmt := SELECT(COUNT(STAR).AS("count")).FROM(stmt.AsTable("c"))

	var count struct{ Count int64 }

	err := Query(ctx, countStmt, &count)
	if err != nil {
		return err
	}

	paging.Size = count.Count

	return nil
}

// PaginateStatement adds limit, offset, and order to stmt.
// The caller runs Count first. An enabled page with size 0 becomes SELECT NULL.
//
//nolint:ireturn // jet SelectStatement is the clause type this module already accepts.
func PaginateStatement(ctx context.Context, stmt SelectStatement, sortableColumns ColumnList) (SelectStatement, error) {
	paging := pagination.Ctx(ctx)

	if !paging.Enabled() {
		return stmt, nil
	}

	if paging.Size == 0 {
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

	return orderStatement(stmt, paging.SortKeys(), sortableColumns)
}

//nolint:ireturn // jet SelectStatement is the clause type this module already accepts.
func orderStatement(
	stmt SelectStatement,
	sortKeys []pagination.SortKey,
	sortableColumns ColumnList,
) (SelectStatement, error) {
	if len(sortKeys) == 0 {
		return stmt, nil
	}

	columns := make(map[string]Column, len(sortableColumns))
	for _, column := range sortableColumns {
		columns[column.Name()] = column
	}

	for _, sort := range sortKeys {
		column, sortable := columns[sort.Key]
		if !sortable {
			return nil, pagination.ErrInvalidSortKey
		}

		switch sort.Order {
		case "", "asc":
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
