package acme

import (
	"crypto/ecdsa"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()

	t.Run("missing key", func(t *testing.T) {
		t.Parallel()

		account, err := NewAccount(t.Context(), filepath.Join(t.TempDir(), "missing.key"), false)
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Nil(t, account.PrivateKey)
	})

	t.Run("creates key", func(t *testing.T) {
		t.Parallel()

		keyPath := filepath.Join(t.TempDir(), "account.key")
		account, err := NewAccount(t.Context(), keyPath, true)
		require.NoError(t, err)
		require.NotNil(t, account.PrivateKey)
		require.True(t, account.TermsOfServiceAgreed)

		loaded, err := NewAccount(t.Context(), keyPath, false)
		require.NoError(t, err)
		require.False(t, loaded.TermsOfServiceAgreed)

		createdKey, createdOK := account.PrivateKey.(*ecdsa.PrivateKey)
		loadedKey, loadedOK := loaded.PrivateKey.(*ecdsa.PrivateKey)

		require.True(t, createdOK)
		require.True(t, loadedOK)
		require.Equal(t, createdKey.PublicKey, loadedKey.PublicKey)
	})

	t.Run("rejects other pem type", func(t *testing.T) {
		t.Parallel()

		keyPath := filepath.Join(t.TempDir(), "account.key")
		encoded := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("not-a-key")})
		require.NoError(t, os.WriteFile(keyPath, encoded, 0o600))

		_, err := NewAccount(t.Context(), keyPath, false)
		require.Error(t, err)
	})
}

func TestCreatePrivateKey(t *testing.T) {
	t.Parallel()

	keyPath := filepath.Join(t.TempDir(), "account.key")
	key, err := createPrivateKey(keyPath)
	require.NoError(t, err)
	require.NotNil(t, key)

	_, err = createPrivateKey(keyPath)
	require.Error(t, err)
}

func TestGetPrivateKey(t *testing.T) {
	t.Parallel()

	keyPath := filepath.Join(t.TempDir(), "account.key")
	_, err := getPrivateKey(keyPath)
	require.ErrorIs(t, err, os.ErrNotExist)

	created, err := createPrivateKey(keyPath)
	require.NoError(t, err)

	loaded, err := getPrivateKey(keyPath)
	require.NoError(t, err)

	createdKey, createdOK := created.(*ecdsa.PrivateKey)
	loadedKey, loadedOK := loaded.(*ecdsa.PrivateKey)

	require.True(t, createdOK)
	require.True(t, loadedOK)
	require.Equal(t, createdKey.PublicKey, loadedKey.PublicKey)
}
