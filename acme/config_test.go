package acme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	require.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Directory)
	assert.NotEmpty(t, cfg.KeyPath)
	assert.Positive(t, cfg.MaxLabels)
	assert.Positive(t, cfg.PollInterval)
	assert.Greater(t, cfg.PollTimeout, cfg.PollInterval)
	assert.NotEmpty(t, cfg.UserAgent)
}
