package crypto

import (
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashCertificate(t *testing.T) {
	t.Parallel()

	raw := []byte("certificate bytes")
	cert := &x509.Certificate{Raw: raw}

	require := require.New(t)
	assert := assert.New(t)

	sum := HashCertificate(cert)
	require.NotEmpty(sum)
	assert.Equal(fmt.Sprintf("%x", sha256.Sum256(raw)), sum)
}
