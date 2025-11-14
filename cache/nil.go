package cache

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
)

var once = sync.Once{}

type Nil[K comparable, V any] struct{}

func (Nil[K, V]) Create(ctx context.Context, k K, v V) {}

func (Nil[K, V]) Delete(ctx context.Context, k K) error {
	return nil
}

func (Nil[K, V]) Read(ctx context.Context, k K) (V, bool) {
	once.Do(func() {
		logger := zerolog.Ctx(ctx)
		logger.Warn().Msg("using nil cache implementation, this is not recommended for production")
	})
	var res V
	return res, false
}

func (Nil[K, V]) Update(ctx context.Context, k K, v V) error {
	return nil
}
