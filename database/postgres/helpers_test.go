package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	err := Exec(t.Context(), nil)
	require.Error(t, err)
}

func TestInsert(t *testing.T) {
	id, err := Insert(t.Context(), nil)
	require.Error(t, err)
	require.Zero(t, id)
}

func TestPaginateStatement(t *testing.T) {
	got, err := PaginateStatement(t.Context(), nil, nil)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestQuery(t *testing.T) {
	err := Query(t.Context(), nil, nil)
	require.Error(t, err)
}
