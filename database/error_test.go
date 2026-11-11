package database

import (
	"testing"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

const (
	mysqlForeignKey uint16 = 1452
	mysqlDuplicate  uint16 = 1062
)

type stubDBError struct{}

func (stubDBError) Error() string { return "untranslated" }

func TestTranslateError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "no rows", err: qrm.ErrNoRows, want: ErrNotFound},
		{name: "mysql foreign key", err: &mysql.MySQLError{Number: mysqlForeignKey}, want: ErrForeignKeyConstrainError},
		{name: "mysql duplicate", err: &mysql.MySQLError{Number: mysqlDuplicate}, want: ErrAlreadyExists},
		{name: "postgres foreign key", err: &pgconn.PgError{Code: "23503"}, want: ErrForeignKeyConstrainError},
		{name: "postgres duplicate", err: &pgconn.PgError{Code: "23505"}, want: ErrAlreadyExists},
		{name: "unknown", err: stubDBError{}, want: stubDBError{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.ErrorIs(t, TranslateError(test.err), test.want)
		})
	}
}
