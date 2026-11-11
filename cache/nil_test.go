package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilCreate(t *testing.T) {
	cache := Nil[string, int]{}

	cache.Create(t.Context(), "a", 1)
}

func TestNilDelete(t *testing.T) {
	cache := Nil[string, int]{}

	err := cache.Delete(t.Context(), "a")
	require.NoError(t, err)
}

func TestNilRead(t *testing.T) {
	cache := Nil[string, int]{}

	got, ok := cache.Read(t.Context(), "a")
	assert.False(t, ok)
	assert.Zero(t, got)

	_, ok = cache.Read(t.Context(), "b")
	assert.False(t, ok)
}

func TestNilUpdate(t *testing.T) {
	cache := Nil[string, int]{}

	err := cache.Update(t.Context(), "a", 1)
	require.NoError(t, err)
}
