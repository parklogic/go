package database

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/postgres"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/parklogic/go/log"
)

func New(ctx context.Context, cfg *Configuration) (context.Context, error) {
	if err := context.Cause(ctx); err != nil {
		return ctx, err
	}

	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return ctx, err
	}

	if cfg.ConnMaxIdleTime != 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
	if cfg.ConnMaxLifetime != 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.MaxIdleConns != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if err := db.PingContext(ctx); err != nil {
		return ctx, err
	}

	switch cfg.Driver {
	case "mysql":
		mysql.SetQueryLogger(log.JetMySQLLogger(cfg.SlowQueryThreshold))
	case "pgx":
		postgres.SetQueryLogger(log.JetPostgresLogger(cfg.SlowQueryThreshold))
	}

	return WithContext(ctx, db), nil
}
