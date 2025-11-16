package log

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a new [slog.Logger] using the provided configuration.
// Returns an error if the configuration is invalid.
func New(cfg *Configuration) (*slog.Logger, error) {
	o := &slog.HandlerOptions{AddSource: cfg.AddSource}

	switch strings.ToLower(cfg.Level) {
	case "debug":
		o.Level = slog.LevelDebug
	case "info":
		o.Level = slog.LevelInfo
	case "warn":
		o.Level = slog.LevelWarn
	case "error":
		o.Level = slog.LevelError
	default:
		return nil, ErrInvalidConfig{Field: "level", Value: cfg.Level}
	}

	var h slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "json":
		h = slog.NewJSONHandler(os.Stderr, o)
	case "logfmt":
		h = slog.NewTextHandler(os.Stderr, o)

	default:
		return nil, ErrInvalidConfig{Field: "format", Value: cfg.Format}
	}

	return slog.New(h), nil
}
