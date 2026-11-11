package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type requiredName struct {
	Name string `validate:"required"`
}

func TestCtx(t *testing.T) {
	t.Parallel()

	validator := Ctx(t.Context())
	require.NotNil(t, validator)
	require.Error(t, validator.StructCtx(t.Context(), requiredName{}))
	require.NoError(t, validator.StructCtx(t.Context(), requiredName{Name: "a"}))
}

func TestWithContext(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator()
	require.NoError(t, err)

	ctx := WithContext(t.Context(), validator)
	err = Ctx(ctx).StructCtx(ctx, hostName{Name: "*.example.com"})
	require.NoError(t, err)
}
