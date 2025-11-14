package database

import (
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAlreadyExists            = errors.New("duplicated key")
	ErrForeignKeyConstrainError = errors.New("foreign key error")
	ErrNotFound                 = errors.New("not found")
)

func TranslateError(err error) error {
	switch {
	case errors.Is(err, qrm.ErrNoRows):
		return ErrNotFound
	}

	if errVal := (&mysql.MySQLError{}); errors.As(err, &errVal) {
		return translateMySQLError(errVal)
	}

	if errVal := (&pgconn.PgError{}); errors.As(err, &errVal) {
		return translatePgError(errVal)
	}

	return err
}

func translateMySQLError(err *mysql.MySQLError) error {
	switch err.Number {
	case 1452:
		return ErrForeignKeyConstrainError
	case 1062:
		return ErrAlreadyExists
	default:
		return err
	}
}

func translatePgError(err *pgconn.PgError) error {
	switch err.Code {
	case "23503":
		return ErrForeignKeyConstrainError
	case "23505":
		return ErrAlreadyExists
	default:
		return err
	}
}
