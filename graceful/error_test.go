package graceful

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorSentinels(t *testing.T) {
	t.Parallel()

	require.NotErrorIs(t, ErrForcedSignal, ErrForcedTimeout)
	require.NotErrorIs(t, ErrForcedTimeout, ErrGracefulSignal)
	require.NotErrorIs(t, ErrGracefulSignal, ErrGracefulExit)
	require.NotErrorIs(t, ErrGracefulExit, ErrForcedSignal)
}
