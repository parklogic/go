package log

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTracedReader(t *testing.T) {
	t.Parallel()

	raw := strings.NewReader("hi")
	assert.Equal(t, raw, NewTracedReader(t.Context(), raw))
}

func TestTracedReaderRead(t *testing.T) {
	t.Parallel()

	var written bytes.Buffer

	logger := zerolog.New(&written).Level(zerolog.TraceLevel)
	reader := NewTracedReader(logger.WithContext(t.Context()), strings.NewReader("hi"))
	got := make([]byte, 2)
	n, err := reader.Read(got)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	assert.Equal(t, "hi", string(got))
	assert.Contains(t, written.String(), "Data received")
}

func TestNewTracedWriter(t *testing.T) {
	t.Parallel()

	var raw bytes.Buffer

	assert.Equal(t, &raw, NewTracedWriter(t.Context(), &raw))
}

func TestTracedWriterWrite(t *testing.T) {
	t.Parallel()

	var written bytes.Buffer

	logger := zerolog.New(&written).Level(zerolog.TraceLevel)

	var raw bytes.Buffer

	writer := NewTracedWriter(logger.WithContext(t.Context()), &raw)
	n, err := writer.Write([]byte("hi"))
	require.NoError(t, err)
	require.Equal(t, 2, n)
	assert.Equal(t, "hi", raw.String())
	assert.Contains(t, written.String(), "Data sent")
}
