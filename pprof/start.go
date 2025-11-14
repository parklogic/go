package pprof

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/rs/zerolog"

	"github.com/parklogic/go/graceful"
)

var Handler *http.ServeMux

// Reset the default server mux, so the pprof endpoints are not exposed.
func init() {
	Handler = http.NewServeMux()
	Handler.Handle("GET /", http.DefaultServeMux)

	http.DefaultServeMux = http.NewServeMux()
}

func Start(ctx context.Context, cfg *Configuration) error {
	if !cfg.Enabled {
		return nil
	}

	logger := zerolog.Ctx(ctx)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("error starting profiling listener on %q: %w", cfg.Address, err)
	}

	server := &http.Server{Addr: cfg.Address, Handler: Handler}
	defer func() {
		if err := server.Close(); err != nil {
			logger.Error().Err(err).Msg("Error closing profiling server")
		}
	}()

	context.AfterFunc(ctx, func() {
		logger.Trace().Str("cause", context.Cause(ctx).Error()).Msg("Closing profiling server")

		shutdownCtx, err := graceful.Ctx(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating profiling shutdown context")
			return
		}

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("Error shutting down profiling server")
			return
		}
	})

	logger.Warn().Stringer("addr", listener.Addr()).Msg("Starting profiling server")
	defer logger.Debug().Msg("Stopped profiling server")
	err = server.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
