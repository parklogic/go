package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ctx, err := New(t.Context())
	require.Error(t, err)
	require.Nil(t, ctx)
}
