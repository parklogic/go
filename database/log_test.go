package database

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type stubStatement struct{}

func (stubStatement) Sql() (string, []any) { return "select 1", nil } //nolint:revive // jet names this method Sql

func (stubStatement) DebugSql() string { return "select 1" } //nolint:revive // jet names this method DebugSql

func TestMySQLLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     zerolog.Level
		info      mysql.QueryInfo
		threshold time.Duration
		want      string
	}{
		{name: "disabled", level: zerolog.Disabled, info: mysql.QueryInfo{}},
		{
			name:      "fast",
			level:     zerolog.WarnLevel,
			threshold: time.Second,
			info:      mysql.QueryInfo{Duration: time.Millisecond},
		},
		{
			name:      "error",
			level:     zerolog.ErrorLevel,
			threshold: time.Second,
			info: mysql.QueryInfo{
				Statement: stubStatement{},
				Err:       stubDBError{},
				Duration:  time.Millisecond,
			},
			want: "Error executing query",
		},
		{
			name:      "ignored rows",
			level:     zerolog.ErrorLevel,
			threshold: time.Second,
			info:      mysql.QueryInfo{Statement: stubStatement{}, Err: qrm.ErrNoRows},
		},
		{
			name:      "slow",
			level:     zerolog.WarnLevel,
			threshold: time.Millisecond,
			info:      mysql.QueryInfo{Statement: stubStatement{}, Duration: time.Second},
			want:      "Slow query",
		},
		{
			name:      "trace",
			level:     zerolog.TraceLevel,
			threshold: time.Second,
			info:      mysql.QueryInfo{Statement: stubStatement{}, Duration: time.Millisecond},
			want:      "Executed query",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assertQueryLog(t, test.level, test.want, func(ctx context.Context) {
				MySQLLogger(test.threshold)(ctx, test.info)
			})
		})
	}
}

func TestPostgresLogger(t *testing.T) {
	t.Parallel()

	assertQueryLog(t, zerolog.Disabled, "", func(ctx context.Context) {
		PostgresLogger(time.Second)(ctx, postgres.QueryInfo{})
	})

	assertQueryLog(t, zerolog.WarnLevel, "", func(ctx context.Context) {
		info := postgres.QueryInfo{Duration: time.Millisecond}
		PostgresLogger(time.Second)(ctx, info)
	})

	assertQueryLog(t, zerolog.ErrorLevel, "Error executing query", func(ctx context.Context) {
		info := postgres.QueryInfo{Statement: stubStatement{}, Err: stubDBError{}, Duration: time.Millisecond}
		PostgresLogger(time.Second)(ctx, info)
	})

	assertQueryLog(t, zerolog.ErrorLevel, "", func(ctx context.Context) {
		info := postgres.QueryInfo{Statement: stubStatement{}, Err: qrm.ErrNoRows}
		PostgresLogger(time.Second)(ctx, info)
	})

	assertQueryLog(t, zerolog.WarnLevel, "Slow query", func(ctx context.Context) {
		info := postgres.QueryInfo{Statement: stubStatement{}, Duration: time.Second}
		PostgresLogger(time.Millisecond)(ctx, info)
	})

	assertQueryLog(t, zerolog.TraceLevel, "Executed query", func(ctx context.Context) {
		info := postgres.QueryInfo{Statement: stubStatement{}, Duration: time.Millisecond}
		PostgresLogger(time.Second)(ctx, info)
	})
}

func assertQueryLog(t *testing.T, level zerolog.Level, want string, run func(context.Context)) {
	t.Helper()

	var written bytes.Buffer

	logger := zerolog.New(&written).Level(level)
	run(logger.WithContext(t.Context()))

	if want == "" {
		assert.Empty(t, written.String())

		return
	}

	assert.Contains(t, written.String(), want)
}
