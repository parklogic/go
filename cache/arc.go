package cache

import (
	"context"

	"github.com/hashicorp/golang-lru/arc/v2"
)

type ARC[K comparable, V any] struct {
	*arc.ARCCache[K, V]
}

func NewARC[K comparable, V any](size int) (*ARC[K, V], error) {
	c, err := arc.NewARC[K, V](size)
	if err != nil {
		return nil, err
	}

	return &ARC[K, V]{c}, nil
}

func (c *ARC[K, V]) Create(ctx context.Context, k K, v V) {
	c.ARCCache.Add(k, v)
}

func (c *ARC[K, V]) Delete(ctx context.Context, k K) error {
	c.ARCCache.Remove(k)

	return nil
}

func (c *ARC[K, V]) Read(ctx context.Context, k K) (V, bool) {
	return c.ARCCache.Get(k)
}

func (c *ARC[K, V]) Update(ctx context.Context, k K, v V) error {
	c.ARCCache.Add(k, v)

	return nil
}
