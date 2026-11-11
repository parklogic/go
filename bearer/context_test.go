package bearer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtx(t *testing.T) {
	_, err := Ctx(t.Context())
	require.Error(t, err)
}

func TestWithContext(t *testing.T) {
	validator := staticValidator{valid: true}
	ctx := WithContext(t.Context(), validator)

	got, err := Ctx(ctx)
	require.NoError(t, err)
	require.Equal(t, validator, got)
}
