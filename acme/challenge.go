package acme

import (
	"context"

	"github.com/mholt/acmez/v3/acme"
)

func InitiateChallenge(ctx context.Context, c acme.Challenge) (challenge acme.Challenge, err error) {
	a, err := Ctx(ctx)
	if err != nil {
		return challenge, err
	}

	return a.Client.InitiateChallenge(ctx, a.Account, c)
}
