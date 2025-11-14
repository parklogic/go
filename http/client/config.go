package client

import (
	"fmt"
	"runtime"
	"time"
)

type Configuration struct {
	ConnectionTimeout     time.Duration `description:"Maximum amount of time waiting for a TCP connection to complete"`
	IdleConnTimeout       time.Duration `description:"Time before closing idle connections"`
	KeepAlive             time.Duration `description:"Interval between keep-alive probes"`
	MaxIdleConnPerHost    int           `description:"Maximum number of idle connections per host"`
	RequestTimeout        time.Duration `description:"Maximum time waiting for requests to complete"`
	ResponseHeaderTimeout time.Duration `description:"Maximum time waiting for response headers"`
	TLSHandshakeTimeout   time.Duration `description:"Maximum time waiting for TLS handshakes to complete"`
	UserAgent             string        `description:"User-Agent header to send with HTTP requests"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		ConnectionTimeout:     30 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		KeepAlive:             30 * time.Second,
		MaxIdleConnPerHost:    runtime.NumCPU(),
		RequestTimeout:        2 * time.Minute,
		ResponseHeaderTimeout: 30 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		UserAgent:             fmt.Sprintf("\"Mozilla/5.0 (compatible; %s)", runtime.Version()),
	}
}
