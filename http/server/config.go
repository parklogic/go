package server

import (
	"time"
)

type Configuration struct {
	Address           string        `description:"Address to listen on"`
	ReadTimeout       time.Duration `description:"Maximum time waiting for the entire request to be read"`
	ReadHeaderTimeout time.Duration `description:"Maximum time waiting for the request headers to be read"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Address:           ":8080",
		ReadTimeout:       0,
		ReadHeaderTimeout: 15 * time.Second,
	}
}
