package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type sampleConfigError struct{}

func (sampleConfigError) Error() string { return "invalid config" }

func (sampleConfigError) Field() string { return "port" }

func (sampleConfigError) Value() string { return "nope" }

func TestErrInvalidConfig(t *testing.T) {
	t.Parallel()

	var err ErrInvalidConfig = sampleConfigError{}

	assert.Equal(t, "port", err.Field())
	assert.Equal(t, "nope", err.Value())
	assert.EqualError(t, err, "invalid config")
}
