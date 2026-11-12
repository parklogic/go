// Package dbctx resolves the statement executor for the database module.
// A transaction in the context wins. Otherwise the function returns the
// database in the context.
package dbctx

import (
	"context"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/parklogic/go/database"
	"github.com/parklogic/go/database/transaction"
)

// DB returns the transaction from ctx, or the database when no transaction is present.
//
//nolint:ireturn // *sql.DB and *sql.Tx share qrm.DB and have no common concrete type.
func DB(ctx context.Context) (qrm.DB, error) {
	if tx := transaction.Ctx(ctx); tx != nil {
		return tx, nil
	}

	// The context error is returned unchanged. Callers used to return database.Ctx's error.
	return database.Ctx(ctx) //nolint:wrapcheck // keep the context error's text
}
