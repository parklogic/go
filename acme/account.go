package acme

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/mholt/acmez/v3/acme"
	"github.com/rs/zerolog"
)

func NewAccount(ctx context.Context, keyPath string, allowCreate bool) (account acme.Account, err error) {
	logger := zerolog.Ctx(ctx)

	key, err := getPrivateKey(keyPath)
	if errors.Is(err, os.ErrNotExist) && allowCreate {
		logger.Warn().Str("path", keyPath).Msg("Account key does not exist, creating new one")

		key, err = createPrivateKey(keyPath)
	}
	if err != nil {
		return account, err
	}

	return acme.Account{
		PrivateKey:           key,
		TermsOfServiceAgreed: allowCreate,
	}, nil
}

func createPrivateKey(p string) (crypto.Signer, error) {
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	bs, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	if err := pem.Encode(f, &pem.Block{Type: "PRIVATE KEY", Bytes: bs}); err != nil {
		return nil, err
	}

	return key, nil
}

func getPrivateKey(p string) (crypto.Signer, error) {
	var key crypto.Signer

	bs, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	// todo: check if block is nil
	block, _ := pem.Decode(bs)

	// todo: support multiple key types
	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("invalid PEM block type: %s", block.Type)
	}

	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := k.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("invalid account key type: %T", k)
	}

	return key, nil
}
