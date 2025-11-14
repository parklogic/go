package acme

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"

	"github.com/mholt/acmez/v3/acme"
)

func NewOrder(ctx context.Context, fqdn string) (order acme.Order, err error) {
	// todo: handle ARI renewals

	a, err := Ctx(ctx)
	if err != nil {
		return order, err
	}

	order = acme.Order{Identifiers: []acme.Identifier{{Type: "dns", Value: fqdn}}}

	return a.Client.NewOrder(ctx, a.Account, order)
}

func GetOrder(ctx context.Context, l string) (order acme.Order, err error) {
	a, err := Ctx(ctx)
	if err != nil {
		return order, err
	}

	order, err = a.Client.GetOrder(ctx, a.Account, acme.Order{Location: l})
	if err != nil {
		return order, err
	}

	return order, nil
}

func FinalizeOrder(ctx context.Context, o acme.Order, pkey crypto.Signer) (order acme.Order, err error) {
	a, err := Ctx(ctx)
	if err != nil {
		return order, err
	}

	// todo: handle IP certificates?

	dnsNames := make([]string, 0, len(o.Identifiers))
	for _, id := range o.Identifiers {
		if id.Type == "dns" {
			dnsNames = append(dnsNames, id.Value)
		}
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		DNSNames: dnsNames,
	}, pkey)
	if err != nil {
		return order, err
	}

	return a.Client.FinalizeOrder(ctx, a.Account, o, csr)
}
