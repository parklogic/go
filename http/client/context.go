package client

import (
	"context"
	"net/http"
)

type contextKey struct{}

func Ctx(ctx context.Context) *http.Client {
	c, ok := ctx.Value(contextKey{}).(*http.Client)
	if !ok {
		return http.DefaultClient
	}

	return c
}

func WithContext(ctx context.Context, c *http.Client) context.Context {
	return context.WithValue(ctx, contextKey{}, c)
}
