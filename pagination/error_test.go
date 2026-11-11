package pagination

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorSentinels(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, ErrInvalidLimit, ErrInvalidPage)
	require.NotErrorIs(t, ErrInvalidLimit, ErrMissingSortKey)
	require.NotErrorIs(t, ErrInvalidPage, ErrInvalidSortKey)
	require.NotErrorIs(t, ErrInvalidSortKey, ErrInvalidSortOrder)
	require.NotErrorIs(t, ErrInvalidSortOrder, ErrMissingSortKey)
}
