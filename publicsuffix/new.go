package publicsuffix

import (
	"context"
)

func New(ctx context.Context, cfg *Configuration) (context.Context, error) {
	if err := context.Cause(ctx); err != nil {
		return ctx, err
	}

	psl, err := NewList(cfg.Path, cfg.CacheSize, cfg.ErrCacheSize)
	if err != nil {
		return ctx, err
	}

	return psl.WithContext(ctx), nil
}
