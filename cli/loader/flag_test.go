package loader

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlagLoad(t *testing.T) {
	t.Parallel()

	cfg := &sampleConfig{}
	loaded, err := (&Flag{}).Load(nil, commandFor(cfg))
	require.NoError(t, err)
	require.False(t, loaded)
	require.Empty(t, cfg.Name)

	cfg = &sampleConfig{}
	loaded, err = (&Flag{}).Load([]string{"--name=from-flag"}, commandFor(cfg))
	require.NoError(t, err)
	require.False(t, loaded)
	require.Equal(t, "from-flag", cfg.Name)

	cfg = &sampleConfig{}
	loaded, err = (&Flag{}).Load([]string{"--=bad"}, commandFor(cfg))
	require.Error(t, err)
	require.True(t, loaded)
}
