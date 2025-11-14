package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
)

type contextKey struct{}

// Ctx ensures a graceful shutdown by monitoring a nested context.
// Returns a detached context, or an error if the required graceful context is missing.
// Automatically cancels the detached context when the graceful context is cancelled.
//
// The detached context will be canceled when the graceful context is cancelled (usually after the main context does not finish after the graceful timeout)
func Ctx(ctx context.Context) (context.Context, error) {
	gracefulCtx, ok := ctx.Value(contextKey{}).(context.Context)
	if !ok {
		return ctx, errors.New("graceful context not found in context")
	}

	detachedCtx, forceStop := context.WithCancelCause(context.WithoutCancel(ctx))

	stopMonitoring := context.AfterFunc(gracefulCtx, func() {
		forceStop(context.Cause(gracefulCtx))
	})

	context.AfterFunc(detachedCtx, func() {
		stopMonitoring()
	})

	return detachedCtx, nil
}

func WithContext(parent context.Context, timeout time.Duration, signals ...os.Signal) (ctx context.Context, cancel context.CancelCauseFunc, forceCancel context.CancelCauseFunc) {
	logger := zerolog.Ctx(parent)

	ctx, cancel = context.WithCancelCause(parent)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, signals...)
		defer signal.Stop(sig)

		select {
		case <-ctx.Done():
			logger.Debug().Err(context.Cause(ctx)).Msg("Main shutdown context cancelled, stopping monitoring signals")
			return

		case s := <-sig:
			logger.Info().Stringer("signal", s).Msg("Received signal, shutting down")
			cancel(ErrGracefulSignal)
		}
	}()

	gracefulCtx, forceCancel := context.WithCancelCause(context.Background())

	stopMonitoring := context.AfterFunc(ctx, func() {
		go func() {
			t := time.AfterFunc(timeout, func() {
				logger.Warn().Msg("Shutdown timeout reached, forcing shutdown")
				forceCancel(ErrForcedTimeout)
			})
			defer t.Stop()

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, signals...)
			defer signal.Stop(sig)

			select {
			case <-gracefulCtx.Done():
				logger.Debug().Err(context.Cause(ctx)).Msg("Graceful shutdown context cancelled, stopping monitoring signals")
				return

			case s := <-sig:
				logger.Warn().Stringer("signal", s).Msg("Received signal, forcing shutdown")
				forceCancel(ErrForcedSignal)
			}
		}()
	})

	context.AfterFunc(gracefulCtx, func() {
		stopMonitoring()
	})

	return context.WithValue(ctx, contextKey{}, gracefulCtx), cancel, forceCancel
}
