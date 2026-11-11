package loader

import (
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

type sampleConfig struct {
	Name string `description:"Sample name" file:"name" json:"name" toml:"name" yaml:"name"`
}

func commandFor(cfg *sampleConfig) *cli.Command {
	return &cli.Command{Name: "demo", Configuration: cfg}
}
