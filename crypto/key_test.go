package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrivateKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		algorithm string
		wantRSA   bool
		wantECDSA bool
		wantErr   bool
	}{
		{name: "rsa", algorithm: "rsa", wantRSA: true},
		{name: "RSA upper", algorithm: "RSA", wantRSA: true},
		{name: "ecdsa", algorithm: "ecdsa", wantECDSA: true},
		{name: "unsupported", algorithm: "dsa", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require := require.New(t)
			assert := assert.New(t)

			key, err := NewPrivateKey(test.algorithm)
			if test.wantErr {
				require.Error(err)
				assert.Nil(key)

				return
			}

			require.NoError(err)
			require.NotNil(key)

			switch {
			case test.wantRSA:
				rsaKey, ok := key.(*rsa.PrivateKey)
				require.True(ok)
				assert.Equal(2048, rsaKey.N.BitLen())
			case test.wantECDSA:
				ecKey, ok := key.(*ecdsa.PrivateKey)
				require.True(ok)
				assert.Equal(elliptic.P256(), ecKey.Curve)
			}
		})
	}
}

func TestHashPublicKey(t *testing.T) {
	t.Parallel()

	t.Run("ecdsa public key", func(t *testing.T) {
		t.Parallel()

		require := require.New(t)
		assert := assert.New(t)

		key, err := NewPrivateKey("ecdsa")
		require.NoError(err)

		sum, err := HashPublicKey(key.Public())
		require.NoError(err)
		assert.Len(sum, 64)
	})

	t.Run("unsupported key", func(t *testing.T) {
		t.Parallel()

		require := require.New(t)

		sum, err := HashPublicKey("not a key")
		require.Error(err)
		require.Empty(sum)
	})
}
