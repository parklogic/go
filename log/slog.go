package log

import (
	"log/slog"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

func GetSlogLevel(l zerolog.Level) slog.Level {
	switch {
	case l <= zerolog.DebugLevel:
		return slog.LevelDebug
	case l == zerolog.InfoLevel:
		return slog.LevelInfo
	case l == zerolog.WarnLevel:
		return slog.LevelWarn
	case l == zerolog.ErrorLevel:
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}

func GetSlogger(l *zerolog.Logger) *slog.Logger {
	return slog.New(slogzerolog.Option{
		Logger: l,
		Level:  slogzerolog.ZeroLogLeveler{Logger: l},
	}.NewZerologHandler())
}

func setSlogDefaultLogger(l *zerolog.Logger) {
	slog.SetDefault(GetSlogger(l))
}
