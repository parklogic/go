package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	cfg := NewConfiguration()
	cfg.ForceColor = true
	got, err := New(ctx, cfg)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, ctx, got)
}
