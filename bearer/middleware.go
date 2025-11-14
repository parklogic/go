package bearer

import (
	"net/http"
	"strings"
)

func Middleware[P any, R any](next func(w http.ResponseWriter, r *http.Request, p P) (s int, res *R, err error)) func(w http.ResponseWriter, r *http.Request, p P) (s int, res *R, err error) {
	return func(w http.ResponseWriter, r *http.Request, p P) (s int, res *R, err error) {
		authHdr := r.Header.Get("Authorization")
		if authHdr == "" {
			return 0, nil, ErrMissingHeader
		}

		authType, secret, _ := strings.Cut(authHdr, " ")

		if authType != "Bearer" {
			return 0, nil, ErrInvalidAuthType
		}

		if secret == "" {
			return 0, nil, ErrMissingToken
		}

		ctx := r.Context()

		validator, err := Ctx(ctx)
		if err != nil {
			return 0, nil, err
		}

		valid, err := validator.Validate(ctx, secret)
		if err != nil {
			return 0, nil, err
		}

		if !valid {
			return 0, nil, ErrInvalidToken
		}

		return next(w, r, p)
	}
}
