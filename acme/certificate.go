package acme

import (
	"context"

	"github.com/mholt/acmez/v3/acme"
)

func GetCertificateChain(ctx context.Context, o acme.Order) (chain acme.Certificate, err error) {
	a, err := Ctx(ctx)
	if err != nil {
		return chain, err
	}

	certChains, err := a.Client.GetCertificateChain(ctx, a.Account, o.Certificate)
	if err != nil {
		return chain, err
	}

	return certChains[0], nil
}
