package chi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter returns a new chi.Mux router with the default middlewares attached.
func NewRouter(cfg *Configuration) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.CleanPath)

	r.Use(LogMiddleware(cfg.SlowResponse))

	if cfg.CompressionLevel > 0 {
		r.Use(middleware.Compress(cfg.CompressionLevel))
	}

	r.Use(RecoverMiddleware)

	return r
}
