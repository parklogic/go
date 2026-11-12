package validation

import (
	"context"
)

type contextKey struct{}

func Ctx(ctx context.Context) Validator {
	v, ok := ctx.Value(contextKey{}).(Validator)
	if !ok {
		return registered
	}

	return v
}

func WithContext(ctx context.Context, v Validator) context.Context {
	return context.WithValue(ctx, contextKey{}, v)
}

type Validator interface {
	StructCtx(ctx context.Context, s any) (err error)
}
