package log

import (
	"context"
	"errors"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/rs/zerolog"
)

func JetMySQLLogger(slowQueryThreshold time.Duration) func(context.Context, mysql.QueryInfo) {
	return func(ctx context.Context, info mysql.QueryInfo) {
		logger := zerolog.Ctx(ctx).With().Dur("took", info.Duration).Int64("rows", info.RowsProcessed).Logger()
		ctx = logger.WithContext(ctx)

		if logger.GetLevel() > zerolog.ErrorLevel {
			return
		}

		if info.Err == nil && info.Duration < slowQueryThreshold && logger.GetLevel() > zerolog.TraceLevel {
			return
		}

		query, args := info.Statement.Sql()
		logger = logger.With().Str("query", query).Interface("args", args).Logger()
		ctx = logger.WithContext(ctx)

		if info.Err != nil && !errors.Is(info.Err, qrm.ErrNoRows) {
			logger.Error().Err(info.Err).Msg("Error executing query")
			return
		}

		if info.Duration > slowQueryThreshold {
			logger.Warn().Msg("Slow query")
		}

		logger.Trace().Msg("Executed query")
	}
}

func JetPostgresLogger(slowQueryThreshold time.Duration) func(context.Context, postgres.QueryInfo) {
	return func(ctx context.Context, info postgres.QueryInfo) {
		logger := zerolog.Ctx(ctx).With().Dur("took", info.Duration).Int64("rows", info.RowsProcessed).Logger()
		ctx = logger.WithContext(ctx)

		if logger.GetLevel() > zerolog.ErrorLevel {
			return
		}

		if info.Err == nil && info.Duration < slowQueryThreshold && logger.GetLevel() > zerolog.TraceLevel {
			return
		}

		query, args := info.Statement.Sql()
		logger = logger.With().Str("query", query).Interface("args", args).Logger()
		ctx = logger.WithContext(ctx)

		if info.Err != nil && !errors.Is(info.Err, qrm.ErrNoRows) {
			logger.Error().Err(info.Err).Msg("Error executing query")
			return
		}

		if info.Duration > slowQueryThreshold {
			logger.Warn().Msg("Slow query")
		}

		logger.Trace().Msg("Executed query")
	}
}
