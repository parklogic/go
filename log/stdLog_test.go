package log

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetStdLogger(t *testing.T) {
	previousFlags := log.Flags()
	previousOutput := log.Writer()

	t.Cleanup(func() {
		log.SetFlags(previousFlags)
		log.SetOutput(previousOutput)
	})

	var written bytes.Buffer

	logger := zerolog.New(&written).Level(zerolog.InfoLevel)
	setStdLogger(&logger)
	log.Print("std-hello")
	assert.Contains(t, written.String(), "std-hello")
	assert.NotEqual(t, os.Stderr, log.Writer())
}
