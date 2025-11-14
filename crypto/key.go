package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"strings"
)

type PrivateKey interface {
	Public() crypto.PublicKey
}

func NewPrivateKey(algorithm string) (PrivateKey, error) {
	switch strings.ToUpper(algorithm) {
	case "RSA":
		return rsa.GenerateKey(rand.Reader, 2048)

	case "ECDSA":
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func HashPublicKey(pub crypto.PublicKey) (string, error) {
	bs, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(bs)), nil
}
