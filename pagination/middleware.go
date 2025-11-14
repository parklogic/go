package pagination

import (
	"net/http"
)

func Middleware[P any, R any](next func(w http.ResponseWriter, r *http.Request, p P) (s int, res *R, err error), cfg ...func(configuration *Configuration)) func(w http.ResponseWriter, r *http.Request, p P) (s int, res *R, err error) {
	config := newConfiguration(cfg...)

	return func(w http.ResponseWriter, r *http.Request, p P) (s int, res *R, err error) {
		pagination, err := FromRequest(config, r)
		if err != nil {
			return 0, nil, err
		}

		ctx := pagination.WithContext(r.Context())

		return next(w, r.WithContext(ctx), p)
	}
}
