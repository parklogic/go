package mysql

import (
	"testing"

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

func TestQuery(t *testing.T) {
	err := Query(t.Context(), nil, nil)
	require.Error(t, err)
}
