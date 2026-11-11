package publicsuffix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewList(t *testing.T) {
	t.Parallel()

	_, err := NewList(filepathMissing(t), 2, 2)
	require.Error(t, err)

	list, err := NewList(writeList(t), 2, 2)
	require.NoError(t, err)

	first, err := list.Parse("www.example.com")
	require.NoError(t, err)
	second, err := list.Parse("www.example.com")
	require.NoError(t, err)
	require.Equal(t, first, second)
	require.Equal(t, "example.com", first.Domain)
	require.Equal(t, "www", first.Subdomain)

	_, err = list.Parse(" ")
	require.Error(t, err)
	_, err = list.Parse(" ")
	require.Error(t, err)
}

func TestParseWithOptions(t *testing.T) {
	t.Parallel()

	list, err := NewList(writeList(t), 0, 0)
	require.NoError(t, err)

	got, err := list.ParseWithOptions("example.com", DefaultFindOptions)
	require.NoError(t, err)
	require.Equal(t, suffixTLD, got.TLD)
}

func TestWithContext(t *testing.T) {
	t.Parallel()

	list, err := NewList(writeList(t), 0, 0)
	require.NoError(t, err)

	loaded, err := Ctx(list.WithContext(t.Context()))
	require.NoError(t, err)
	require.Equal(t, list.psl, loaded.psl)
}

func filepathMissing(t *testing.T) string {
	t.Helper()

	return writeList(t) + ".missing"
}
