package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog"
)

func Rollback(ctx context.Context) {
	tx := Ctx(ctx)

	if tx == nil {
		return
	}

	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		logger := zerolog.Ctx(ctx)
		logger.Error().Err(err).Msg("failed to rollback database transaction")
	}
}
