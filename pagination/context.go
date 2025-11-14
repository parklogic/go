package pagination

import (
	"context"
)

type contextKey struct{}

func Ctx(ctx context.Context) *Pagination {
	p, ok := ctx.Value(contextKey{}).(*Pagination)
	if ok {
		return p
	}

	return &Pagination{}
}
