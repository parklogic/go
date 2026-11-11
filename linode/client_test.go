package linode

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestClientLoggerErrorf(t *testing.T) {
	t.Parallel()

	var written bytes.Buffer

	logger := zerolog.New(&written).Level(zerolog.TraceLevel)
	clientLogger{l: &logger}.Errorf("failed %s", "peer")
	assert.Contains(t, written.String(), "failed peer")
}

func TestClientLoggerWarnf(t *testing.T) {
	t.Parallel()

	var written bytes.Buffer

	logger := zerolog.New(&written).Level(zerolog.TraceLevel)
	clientLogger{l: &logger}.Warnf("slow %s", "peer")
	assert.Contains(t, written.String(), "slow peer")
}

func TestClientLoggerDebugf(t *testing.T) {
	t.Parallel()

	var written bytes.Buffer

	logger := zerolog.New(&written).Level(zerolog.TraceLevel)
	clientLogger{l: &logger}.Debugf("debug %s", "peer")
	assert.Contains(t, written.String(), "debug peer")
}
