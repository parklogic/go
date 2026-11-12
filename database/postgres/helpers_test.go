package postgres

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/parklogic/go/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	err := Exec(t.Context(), nil)
	require.Error(t, err)
}

func TestInsert(t *testing.T) {
	err := Insert(t.Context(), nil)
	require.Error(t, err)
}

func TestCountDisabled(t *testing.T) {
	err := Count(t.Context(), nil)
	require.NoError(t, err)
}

func TestPaginateStatement(t *testing.T) {
	got, err := PaginateStatement(t.Context(), nil, nil)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestPaginateStatementEmptyPage(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/items?limit=10", nil)
	page, err := pagination.FromRequest(pagination.Configuration{
		LimitParameter: "limit",
		PageParameter:  "page",
		SortParameter:  "sort",
		DefaultSort:    "id",
		DefaultLimit:   25,
		MaxLimit:       1000,
	}, req)
	require.NoError(t, err)
	require.True(t, page.Enabled())

	got, err := PaginateStatement(page.WithContext(t.Context()), nil, nil)
	require.NoError(t, err)

	query, args := got.Sql()
	require.Empty(t, args)
	require.Equal(t, "\nSELECT NULL;\n", query)
}

func TestQuery(t *testing.T) {
	err := Query(t.Context(), nil, nil)
	require.Error(t, err)
}
