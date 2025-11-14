package database

import (
	"runtime"
	"time"
)

type Configuration struct {
	Driver string `description:"Database driver"`
	DSN    string `description:"Database connection string"`

	SlowQueryThreshold time.Duration `description:"Threshold for logging slow queries"`

	ConnMaxIdleTime time.Duration `description:"Maximum amount of time a connection may be idle before being closed"`
	ConnMaxLifetime time.Duration `description:"Maximum amount of time a connection may be reused"`
	MaxIdleConns    int           `description:"Maximum number connections in the idle connection pool"`
	MaxOpenConns    int           `description:"Maximum number of open connections to the database"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Driver: "pgx",
		DSN:    "postgresql://postgres:postgres@127.0.0.1/autossl?sslmode=disable",

		SlowQueryThreshold: 500 * time.Millisecond,

		ConnMaxLifetime: 5 * time.Minute,
		MaxIdleConns:    runtime.NumCPU(),
		MaxOpenConns:    runtime.NumCPU() * 2,
	}
}
