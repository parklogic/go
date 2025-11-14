package transaction

import (
	"context"

	"github.com/parklogic/go/database"
	"github.com/parklogic/go/graceful"
)

func New(ctx context.Context) (context.Context, error) {
	db, err := database.Ctx(ctx)
	if err != nil {
		return nil, err
	}

	txCtx, err := graceful.Ctx(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(txCtx, nil)
	if err != nil {
		return nil, err
	}

	return WithContext(ctx, tx), nil
}
