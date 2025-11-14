package transaction

import (
	"context"
)

func Commit(ctx context.Context) error {
	tx := Ctx(ctx)

	if tx == nil {
		return nil
	}

	return tx.Commit()
}
