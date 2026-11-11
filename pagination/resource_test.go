package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResource(t *testing.T) {
	t.Parallel()

	page := Pagination{enabled: true, limit: 2, page: 2, nextPage: "/items?page=3", Size: 2}
	ctx := page.WithContext(t.Context())

	got := NewResource(ctx, []string{"a", "b"})
	require.Equal(t, int64(2), got.Page)
	assert.Equal(t, []string{"a", "b"}, got.Results)
	assert.Equal(t, "/items?page=3", got.Next)
	assert.Equal(t, int64(2), got.Size)
}
