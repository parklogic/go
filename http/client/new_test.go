package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	got := Ctx(New(t.Context(), cfg))
	require.Equal(t, cfg.RequestTimeout, got.Timeout)
	require.NotNil(t, got.Transport)
}
