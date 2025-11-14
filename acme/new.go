package acme

import (
	"context"
	"errors"

	"github.com/mholt/acmez/v3/acme"
	"github.com/rs/zerolog"
)

func New(ctx context.Context, cfg *Configuration) (context.Context, error) {
	if err := context.Cause(ctx); err != nil {
		return ctx, err
	}

	logger := zerolog.Ctx(ctx)

	c := NewClient(cfg)

	a, err := NewAccount(ctx, cfg.KeyPath, cfg.CreateNewAccount)
	if err != nil {
		return nil, err
	}

	a, err = c.GetAccount(ctx, a)
	if p := (acme.Problem{}); errors.As(err, &p) {
		if p.Type == acme.ProblemTypeAccountDoesNotExist {
			logger.Warn().Msg("Account does not exist, creating new one")
			a, err = c.NewAccount(ctx, a)
		}
	}
	if err != nil {
		return nil, err
	}

	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("acme_account", a.Location)
	})

	return contextStruct{
		Account: a,
		Client:  c,
	}.WithContext(ctx), nil
}
