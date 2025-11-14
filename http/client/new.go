package client

import (
	"context"
)

func New(ctx context.Context, cfg *Configuration) context.Context {
	return WithContext(ctx, NewClient(cfg))
}
