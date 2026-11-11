package publicsuffix

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	got, err := New(ctx, NewConfiguration())
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, ctx, got)

	cfg := NewConfiguration()
	cfg.Path = writeList(t)
	loaded, err := New(t.Context(), cfg)
	require.NoError(t, err)

	list, err := Ctx(loaded)
	require.NoError(t, err)

	parsed, err := list.Parse("example.com")
	require.NoError(t, err)
	require.Equal(t, "example.com", parsed.Domain)
}

func writeList(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "public_suffix_list.dat")
	body := "// ===BEGIN ICANN DOMAINS===\ncom\n// ===END ICANN DOMAINS===\n"
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))

	return path
}
