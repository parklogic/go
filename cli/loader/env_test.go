package loader

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvLoad(t *testing.T) {
	loader := &Env{prefix: "PLLOADERTEST_"}
	cfg := &sampleConfig{}
	loaded, err := loader.Load(nil, commandFor(cfg))
	require.NoError(t, err)
	require.False(t, loaded)

	t.Setenv("PLLOADERTEST_NAME", "from-env")

	cfg = &sampleConfig{}
	loaded, err = loader.Load(nil, commandFor(cfg))
	require.NoError(t, err)
	require.False(t, loaded)
	require.Equal(t, "from-env", cfg.Name)
}
