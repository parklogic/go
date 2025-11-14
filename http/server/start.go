package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/parklogic/go/graceful"
	"github.com/parklogic/go/log"
)

func Start(ctx context.Context, cfg *Configuration, handler http.Handler) error {
	logger := zerolog.Ctx(ctx)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("error starting HTTP server listener on %q: %w", cfg.Address, err)
	}

	server := &http.Server{
		Addr:              cfg.Address,
		BaseContext:       log.HTTPBaseContext(ctx),
		ConnContext:       log.HTTPConnContext,
		ConnState:         log.HTTPConnState(ctx),
		ErrorLog:          log.HTTPLogger(ctx),
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
	}
	defer func() {
		if err := server.Close(); err != nil {
			logger.Error().Err(err).Msg("Error closing HTTP server")
		}
	}()

	context.AfterFunc(ctx, shutdown(ctx, server))

	logger.Info().Stringer("addr", listener.Addr()).Msg("Starting HTTP server")
	defer logger.Debug().Msg("Stopped HTTP server")

	if err = server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func shutdown(ctx context.Context, s *http.Server) func() {
	logger := zerolog.Ctx(ctx)

	return func() {
		logger.Trace().Str("cause", context.Cause(ctx).Error()).Msg("Closing HTTP server")

		shutdownCtx, err := graceful.Ctx(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating HTTP server shutdown context")
			return
		}

		if err := s.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("Error shutting down HTTP server")
			return
		}
	}
}
