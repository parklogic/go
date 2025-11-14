package api

import (
	myChi "github.com/parklogic/go/chi"
)

// NewRouter returns a new chi router
func NewRouter(cfg *myChi.Configuration) *chi.Mux {
	r := myChi.NewRouter(cfg)

	r.MethodNotAllowed(Adapter(HandlerMethodNotAllowed))

	r.NotFound(Adapter(HandlerNotFound))

	// Health checks
	r.Get("/livez", Adapter(HandlerLive))

	return r
}
