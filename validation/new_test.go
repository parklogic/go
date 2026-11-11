package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ctx, err := New(t.Context())
	require.NoError(t, err)

	err = Ctx(ctx).StructCtx(ctx, hostName{Name: "example.com"})
	require.NoError(t, err)
}
