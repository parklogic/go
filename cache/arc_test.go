package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewARC(t *testing.T) {
	t.Parallel()

	t.Run("positive size", func(t *testing.T) {
		t.Parallel()

		require := require.New(t)

		cache, err := NewARC[string, int](2)
		require.NoError(err)
		require.NotNil(cache)
	})

	t.Run("zero size", func(t *testing.T) {
		t.Parallel()

		require := require.New(t)

		cache, err := NewARC[string, int](0)
		require.Error(err)
		require.Nil(cache)
	})
}

func TestARCCreate(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	cache, err := NewARC[string, string](2)
	require.NoError(t, err)

	cache.Create(ctx, "a", "one")

	got, ok := cache.Read(ctx, "a")
	assert.True(t, ok)
	assert.Equal(t, "one", got)
}

func TestARCDelete(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	cache, err := NewARC[string, string](2)
	require.NoError(t, err)

	cache.Create(ctx, "a", "one")

	err = cache.Delete(ctx, "a")
	require.NoError(t, err)

	_, ok := cache.Read(ctx, "a")
	assert.False(t, ok)
}

func TestARCRead(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	cache, err := NewARC[string, string](2)
	require.NoError(t, err)

	_, ok := cache.Read(ctx, "missing")
	assert.False(t, ok)
}

func TestARCUpdate(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	cache, err := NewARC[string, string](2)
	require.NoError(t, err)

	cache.Create(ctx, "a", "one")

	err = cache.Update(ctx, "a", "two")
	require.NoError(t, err)

	got, ok := cache.Read(ctx, "a")
	assert.True(t, ok)
	assert.Equal(t, "two", got)
}
