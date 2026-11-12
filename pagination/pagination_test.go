package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		target  string
		wantErr error
		limit   int64
		page    int64
	}{
		{name: "defaults", target: "/items", limit: 25, page: 1},
		{name: "explicit", target: "/items?limit=10&page=2&sort=name,asc", limit: 10, page: 2},
		{name: "bad limit", target: "/items?limit=5000", wantErr: ErrInvalidLimit},
		{name: "bad page", target: "/items?page=0", wantErr: ErrInvalidPage},
		{name: "missing sort key", target: "/items?sort=,asc", wantErr: ErrMissingSortKey},
		{name: "bad sort order", target: "/items?sort=name,sideways", wantErr: ErrInvalidSortOrder},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require := require.New(t)
			assert := assert.New(t)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, test.target, nil)

			got, err := FromRequest(newConfiguration(), req)
			if test.wantErr != nil {
				require.ErrorIs(err, test.wantErr)

				return
			}

			require.NoError(err)
			assert.True(got.Enabled())
			assert.Equal(test.limit, got.Limit())
			assert.Equal(test.page, got.Page())
			assert.NotEmpty(got.SortKeys())
			assert.Equal(test.limit*(test.page-1), got.Offset())

			got.Size = got.Limit()
			assert.Contains(got.NextPage(), "page=")
		})
	}
}

func TestEnabled(t *testing.T) {
	t.Parallel()

	assert.False(t, Pagination{}.Enabled())
}

func TestLimit(t *testing.T) {
	t.Parallel()

	assert.Equal(t, int64(10), Pagination{limit: 10}.Limit())
}

func TestNextPage(t *testing.T) {
	t.Parallel()

	short := Pagination{limit: 10, nextPage: "/items?page=2", Size: 3}
	assert.Empty(t, short.NextPage())

	full := Pagination{limit: 10, nextPage: "/items?page=2", Size: 10}
	assert.Equal(t, "/items?page=2", full.NextPage())
}

func TestOffset(t *testing.T) {
	t.Parallel()

	assert.Equal(t, int64(0), Pagination{limit: 10, page: 1}.Offset())
	assert.Equal(t, int64(20), Pagination{limit: 10, page: 3}.Offset())
}

func TestPage(t *testing.T) {
	t.Parallel()

	assert.Equal(t, int64(1), Pagination{}.Page())
	assert.Equal(t, int64(4), Pagination{page: 4}.Page())
}

func TestSortKeys(t *testing.T) {
	t.Parallel()

	assert.Empty(t, Pagination{}.SortKeys())

	keys := []SortKey{{Key: "name", Order: "asc"}}
	assert.Equal(t, keys, Pagination{sortKeys: keys}.SortKeys())
}

func TestWithContext(t *testing.T) {
	t.Parallel()

	page := Pagination{enabled: true, limit: 10, page: 2}
	got := Ctx(page.WithContext(t.Context()))
	require.NotNil(t, got)
	assert.True(t, got.Enabled())
	assert.Equal(t, int64(10), got.Limit())
}
