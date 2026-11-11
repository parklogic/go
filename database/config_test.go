package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfiguration(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	require.NotNil(t, cfg)
	assert.Equal(t, "pgx", cfg.Driver)
	assert.NotEmpty(t, cfg.DSN)
	assert.Positive(t, cfg.SlowQueryThreshold)
	assert.Positive(t, cfg.MaxOpenConns)
	assert.GreaterOrEqual(t, cfg.MaxOpenConns, cfg.MaxIdleConns)
}
