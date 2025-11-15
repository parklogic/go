package chi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func LogMiddleware(slowThreshold time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logger := zerolog.Ctx(ctx).With().Logger()
			ctx = logger.WithContext(ctx)
			r = r.WithContext(ctx)

			if logger.GetLevel() > zerolog.ErrorLevel {
				next.ServeHTTP(w, r)
				return
			}

			logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.
					Str("host", r.Host).
					Str("method", r.Method).
					Str("referer", r.Referer()).
					Str("request_id", middleware.GetReqID(ctx)).
					Str("source", r.RemoteAddr).
					Str("user_agent", r.UserAgent()).
					Stringer("url", r.URL)
			})

			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				elapsed := time.Since(start)

				logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.
						Dur("elapsed", elapsed).
						Int("sent", ww.BytesWritten()).
						Int("status", ww.Status())
				})

				switch {
				case ww.Status() >= http.StatusInternalServerError:
					logger.Error().Stack().Msg("Request failed")

				case elapsed > slowThreshold:
					logger.Warn().Msg("Slow request")

				case ww.Status() == http.StatusNotFound:
					logger.Trace().Msg("Bad request")

				case ww.Status() >= http.StatusBadRequest:
					logger.Info().Msg("Bad request")

				default:
					logger.Trace().Msg("Request completed")
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				logger := zerolog.Ctx(r.Context())

				logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
					if err, ok := rvr.(error); ok {
						return c.Stack().Err(errors.WithStack(err))
					}

					if s, ok := rvr.(string); ok {
						return c.Stack().Err(errors.WithStack(errors.New(s)))
					}

					return c
				})

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
