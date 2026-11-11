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
	assert.True(t, cfg.AddSource)
	assert.Equal(t, "json", cfg.Format)
	assert.Equal(t, "info", cfg.Level)
}

func TestErrInvalidConfig(t *testing.T) {
	t.Parallel()

	err := errInvalidConfig{field: "level", value: "nope"}
	assert.Equal(t, "level", err.Field())
	assert.Equal(t, "nope", err.Value())
	assert.EqualError(t, err, "invalid level: nope")
}
