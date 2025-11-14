package bearer

import (
	"context"
	"fmt"
)

type contextKey struct{}

func Ctx(ctx context.Context) (Validator, error) {
	validator, ok := ctx.Value(contextKey{}).(Validator)
	if !ok {
		return nil, fmt.Errorf("token validator not found in context")
	}

	return validator, nil
}

func WithContext(ctx context.Context, validator Validator) context.Context {
	return context.WithValue(ctx, contextKey{}, validator)
}

type Validator interface {
	Validate(ctx context.Context, secret string) (bool, error)
}
