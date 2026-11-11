package pagination

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtx(t *testing.T) {
	t.Parallel()

	got := Ctx(t.Context())
	require.NotNil(t, got)
	require.False(t, got.Enabled())
}
