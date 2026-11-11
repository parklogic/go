package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDefaultSize(t *testing.T) {
	t.Parallel()

	cfg := newConfiguration(WithDefaultSize(50))
	assert.Equal(t, int64(50), cfg.DefaultLimit)
}

func TestWithDefaultSort(t *testing.T) {
	t.Parallel()

	cfg := newConfiguration(WithDefaultSort("name,desc"))
	assert.Equal(t, "name,desc", cfg.DefaultSort)
}

func TestWithMaxSize(t *testing.T) {
	t.Parallel()

	cfg := newConfiguration(WithMaxSize(100))
	assert.Equal(t, int64(100), cfg.MaxLimit)
}

func TestWithMinSize(t *testing.T) {
	t.Parallel()

	cfg := newConfiguration(WithMinSize(1))
	assert.Equal(t, int64(1), cfg.MinLimit)
}

func TestWithOptionalPagination(t *testing.T) {
	t.Parallel()

	cfg := newConfiguration(WithOptionalPagination())
	assert.Equal(t, int64(-1), cfg.MinLimit)
}
