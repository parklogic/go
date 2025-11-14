package client

import (
	"net/http"
)

func NewClient(cfg *Configuration) *http.Client {
	return &http.Client{
		Timeout:   cfg.RequestTimeout,
		Transport: newTransport(cfg),
	}
}
