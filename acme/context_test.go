package acme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	t.Parallel()

	stored := contextStruct{}
	got, err := Ctx(stored.WithContext(t.Context()))
	require.NoError(t, err)
	assert.Equal(t, stored, got)
}

func TestCtx(t *testing.T) {
	t.Parallel()

	_, err := Ctx(t.Context())
	require.EqualError(t, err, "ACME client not found in context")
}
