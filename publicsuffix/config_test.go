package publicsuffix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	require.NotNil(t, cfg)
	assert.Positive(t, cfg.CacheSize)
	assert.Positive(t, cfg.ErrCacheSize)
	assert.Equal(t, DefaultListPath, cfg.Path)
}
