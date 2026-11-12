package loader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traefik/paerser/cli"
)

func TestNew(t *testing.T) {
	t.Parallel()

	loaders := New("demo-app")
	require.Len(t, loaders, 3)

	_, fileOK := loaders[0].(*File)
	_, envOK := loaders[1].(*Env)
	_, flagOK := loaders[2].(*Flag)

	require.True(t, fileOK)
	require.True(t, envOK)
	require.True(t, flagOK)
}

func TestCleanEnvPrefix(t *testing.T) {
	t.Parallel()

	require.Equal(t, "DEMOAPP_", cleanEnvPrefix("demo-app"))
}

func TestLoadOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "demo.yaml")
	require.NoError(t, os.WriteFile(path, []byte("name: from-file\ncity: from-file\n"), 0o600))
	t.Setenv("DEMOAPP_NAME", "from-env")

	cfg := &sampleConfig{}
	cmd := commandFor(cfg)
	args := []string{"--config=" + path, "--name=from-flag"}
	wantName := []string{"from-file", "from-env", "from-flag"}

	for i, loader := range New("demo-app") {
		loaded, err := loader.Load(args, cmd)
		require.NoError(t, err)
		require.False(t, loaded)
		require.Equal(t, wantName[i], cfg.Name)
		require.Equal(t, "from-file", cfg.City)
	}
}

type sampleConfig struct {
	Name string `description:"Sample name" file:"name" json:"name" toml:"name" yaml:"name"`
	City string `description:"Sample city" file:"city" json:"city" toml:"city" yaml:"city"`
}

func commandFor(cfg *sampleConfig) *cli.Command {
	return &cli.Command{Name: "demo", Configuration: cfg}
}
