package loader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileLoad(t *testing.T) {
	loader := &File{envPrefix: "PLLOADERFILE_"}
	cfg := &sampleConfig{}
	missing := filepath.Join(t.TempDir(), "missing.yaml")
	args := []string{"--config=" + missing}

	loaded, err := loader.Load(args, commandFor(cfg))
	require.Error(t, err)
	require.True(t, loaded)

	dir := t.TempDir()
	configPath := filepath.Join(dir, "demo.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte("name: from-file\n"), 0o600))

	cfg = &sampleConfig{}
	loaded, err = loader.Load([]string{"--config=" + configPath}, commandFor(cfg))
	require.NoError(t, err)
	require.False(t, loaded)
	require.Equal(t, "from-file", cfg.Name)

	t.Setenv("PLLOADERFILE_CONFIG", configPath)

	cfg = &sampleConfig{}
	loaded, err = loader.Load(nil, commandFor(cfg))
	require.NoError(t, err)
	require.False(t, loaded)
	require.Equal(t, "from-file", cfg.Name)
}
