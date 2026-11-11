package pprof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	cfg := NewConfiguration()
	cfg.Enabled = false

	err := Start(t.Context(), cfg)
	require.NoError(t, err)
}
