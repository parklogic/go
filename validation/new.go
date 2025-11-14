package validation

import (
	"context"
)

func New(ctx context.Context) (context.Context, error) {
	v, err := NewValidator()
	if err != nil {
		return ctx, err
	}

	return WithContext(ctx, v), nil
}
