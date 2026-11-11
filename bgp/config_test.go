package bgp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	require.NotNil(t, cfg)
	assert.Equal(t, ":50051", cfg.Listen)
	assert.Positive(t, cfg.RecvMaxMsgSize)
	assert.Equal(t, cfg.RecvMaxMsgSize, cfg.SendMaxMsgSize)
}
