package log

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSlogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level zerolog.Level
		want  slog.Level
	}{
		{name: "trace", level: zerolog.TraceLevel, want: slog.LevelDebug},
		{name: "debug", level: zerolog.DebugLevel, want: slog.LevelDebug},
		{name: "info", level: zerolog.InfoLevel, want: slog.LevelInfo},
		{name: "warn", level: zerolog.WarnLevel, want: slog.LevelWarn},
		{name: "error", level: zerolog.ErrorLevel, want: slog.LevelError},
		{name: "disabled", level: zerolog.Disabled, want: slog.LevelInfo},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, GetSlogLevel(test.level))
		})
	}
}

func TestGetSlogger(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(bytes.NewBuffer(nil))
	require.NotNil(t, GetSlogger(&logger))
}

func TestSetSlogDefaultLogger(t *testing.T) {
	previous := slog.Default()

	t.Cleanup(func() { slog.SetDefault(previous) })

	logger := zerolog.New(bytes.NewBuffer(nil))
	setSlogDefaultLogger(&logger)
	require.NotEqual(t, previous, slog.Default())
}
