package publicsuffix

import (
	"context"
	"errors"
)

type contextKey struct{}

func Ctx(ctx context.Context) (List, error) {
	var l List

	l, ok := ctx.Value(contextKey{}).(List)
	if !ok {
		return l, errors.New("public suffix list not found in context")
	}

	return l, nil
}
