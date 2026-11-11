package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/parklogic/go/graceful"
	"github.com/stretchr/testify/require"
)

func TestShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	shutdown(ctx, &http.Server{ReadHeaderTimeout: time.Second})()

	parent, stop, force := graceful.WithContext(t.Context(), time.Hour)
	t.Cleanup(func() {
		stop(context.Canceled)
		force(context.Canceled)
	})

	stop(graceful.ErrGracefulExit)
	shutdown(parent, &http.Server{ReadHeaderTimeout: time.Second})()
	require.ErrorIs(t, context.Cause(parent), graceful.ErrGracefulExit)
}
