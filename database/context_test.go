package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtx(t *testing.T) {
	_, err := Ctx(t.Context())
	require.Error(t, err)
}

func TestWithContext(t *testing.T) {
	opened, err := sql.Open("pgx", "postgres://postgres@127.0.0.1/postgres?sslmode=disable")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, opened.Close())
	})

	got, err := Ctx(WithContext(t.Context(), opened))
	require.NoError(t, err)
	require.Equal(t, opened, got)
}
