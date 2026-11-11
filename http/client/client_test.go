package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	got := NewClient(cfg)
	require.Equal(t, cfg.RequestTimeout, got.Timeout)
	require.NotNil(t, got.Transport)
}
