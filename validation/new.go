package validation

import (
	"context"
)

func New(ctx context.Context) (context.Context, error) {
	return WithContext(ctx, registered), nil
}
