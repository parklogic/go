package graceful

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCtx(t *testing.T) {
	_, err := Ctx(t.Context())
	require.EqualError(t, err, "graceful context not found in context")
}

func TestWithContext(t *testing.T) {
	ctx, cancel, forceCancel := WithContext(t.Context(), time.Hour)
	t.Cleanup(func() {
		cancel(ErrGracefulExit)
		forceCancel(ErrForcedTimeout)
	})

	require.NoError(t, ctx.Err())

	cancel(ErrGracefulExit)
	require.ErrorIs(t, context.Cause(ctx), ErrGracefulExit)

	detached, err := Ctx(ctx)
	require.NoError(t, err)

	forceCancel(ErrForcedTimeout)
	require.Eventually(t, func() bool {
		return detached.Err() != nil
	}, time.Second, time.Millisecond)
	require.ErrorIs(t, context.Cause(detached), ErrForcedTimeout)
}
