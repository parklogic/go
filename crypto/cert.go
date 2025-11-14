package crypto

import (
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

func HashCertificate(cert *x509.Certificate) string {
	return fmt.Sprintf("%x", sha256.Sum256(cert.Raw))
}
