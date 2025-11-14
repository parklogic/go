package acme

import (
	"context"

	"github.com/mholt/acmez/v3/acme"
)

func GetAuthorization(ctx context.Context, l string) (authz acme.Authorization, err error) {
	a, err := Ctx(ctx)
	if err != nil {
		return authz, err
	}

	authz, err = a.Client.GetAuthorization(ctx, a.Account, l)
	if err != nil {
		return authz, err
	}

	return authz, nil
}

func PollAuthorization(ctx context.Context, a acme.Authorization) (authz acme.Authorization, err error) {
	s, err := Ctx(ctx)
	if err != nil {
		return authz, err
	}

	return s.Client.PollAuthorization(ctx, s.Account, a)
}
