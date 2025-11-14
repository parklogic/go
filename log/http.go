package log

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/parklogic/go/graceful"
)

var started sync.Map

func HTTPBaseContext(ctx context.Context) func(net.Listener) context.Context {
	return func(l net.Listener) context.Context {
		logger := zerolog.Ctx(ctx).With().Stringer("addr", l.Addr()).Logger()
		ctx = logger.WithContext(ctx)

		if ctx, err := graceful.Ctx(ctx); err == nil {
			return ctx
		}

		return ctx
	}
}

func HTTPConnContext(ctx context.Context, c net.Conn) context.Context {
	logger := zerolog.Ctx(ctx).With().Stringer("local", c.LocalAddr()).Stringer("remote", c.RemoteAddr()).Logger()
	ctx = logger.WithContext(ctx)

	return ctx
}

func HTTPConnState(ctx context.Context) func(net.Conn, http.ConnState) {
	if zerolog.Ctx(ctx).GetLevel() > zerolog.TraceLevel {
		return nil
	}

	logger := zerolog.Ctx(ctx)

	return func(c net.Conn, cs http.ConnState) {

		switch cs {
		case http.StateNew:
			logger.Trace().Msg("New connection")
			started.Store(c, time.Now())

		case http.StateClosed:
			value, loaded := started.LoadAndDelete(c)
			if loaded {
				if t, ok := value.(time.Time); ok {
					logger.Trace().Dur("duration", time.Since(t)).Msg("Connection closed")
					return
				}
			}

			logger.Trace().Msg("Connection closed")

		default:
		}
	}
}

func HTTPLogger(ctx context.Context) *log.Logger {
	return log.New(zerolog.Ctx(ctx), "", 0)
}
