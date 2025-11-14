package acme

import (
	"context"
	"errors"

	"github.com/mholt/acmez/v3/acme"
)

type contextKey struct{}

type contextStruct struct {
	Account acme.Account
	Client  *acme.Client
}

func (c contextStruct) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, c)
}

func Ctx(ctx context.Context) (contextStruct, error) {
	var c contextStruct

	c, ok := ctx.Value(contextKey{}).(contextStruct)
	if !ok {
		return c, errors.New("ACME client not found in context")
	}

	return c, nil
}
