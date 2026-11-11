package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	require.NotNil(t, cfg)
	assert.Positive(t, cfg.ConnectionTimeout)
	assert.Positive(t, cfg.IdleConnTimeout)
	assert.Positive(t, cfg.KeepAlive)
	assert.Positive(t, cfg.MaxIdleConnPerHost)
	assert.Positive(t, cfg.RequestTimeout)
	assert.Positive(t, cfg.ResponseHeaderTimeout)
	assert.Positive(t, cfg.TLSHandshakeTimeout)
	assert.NotEmpty(t, cfg.UserAgent)
}
