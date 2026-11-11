package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	require.NotNil(t, cfg)
	assert.Equal(t, "info", cfg.Level)
}
