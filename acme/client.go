package acme

import (
	"log/slog"

	"github.com/mholt/acmez/v3/acme"
)

func NewClient(cfg *Configuration) *acme.Client {
	return &acme.Client{
		Directory:    cfg.Directory,
		Logger:       slog.Default(),
		PollInterval: cfg.PollInterval,
		PollTimeout:  cfg.PollTimeout,
		UserAgent:    cfg.UserAgent,
	}
}
