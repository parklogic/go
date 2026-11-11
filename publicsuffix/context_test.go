package publicsuffix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtx(t *testing.T) {
	t.Parallel()

	_, err := Ctx(t.Context())
	require.EqualError(t, err, "public suffix list not found in context")
}
