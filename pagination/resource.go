package pagination

import (
	"context"
)

// Resource represents a paginated list of resources
type Resource[T any] struct {
	// URL to the next page
	Next string `json:"next,omitempty" example:"http://example.com/resource?page=2&limit=25"`
	// Current page
	Page int64 `json:"page" example:"1"`
	// List of resources
	Results []T `json:"results"`
	// Total number of resources
	Size int64 `json:"size" example:"100"`
}

func NewResource[T any](ctx context.Context, res []T) Resource[T] {
	paging := Ctx(ctx)

	return Resource[T]{
		Next:    paging.NextPage(),
		Page:    paging.Page(),
		Results: res,
		Size:    paging.Size,
	}
}
