package log

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func New(ctx context.Context, cfg *Configuration) (context.Context, error) {
	if err := context.Cause(ctx); err != nil {
		return ctx, err
	}

	if cfg.ForceColor || isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd()) {
		zlog.Logger = zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = os.Stderr
			w.TimeFormat = time.StampMilli
		}))
	}

	zlog.Logger = zlog.Logger.With().Caller().Timestamp().Logger()

	setStdLogger(&zlog.Logger)
	setSlogDefaultLogger(&zlog.Logger)

	zerolog.DefaultContextLogger = &zlog.Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if cfg.RemovePrefix != "" {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return strings.TrimPrefix(file, cfg.RemovePrefix) + ":" + strconv.Itoa(line)
		}
	}

	if cfg.ShortenPath {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return filepath.Base(file) + ":" + strconv.Itoa(line)
		}
	}

	logLevel, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return ctx, fmt.Errorf("invalid log level: %w", err)
	}

	zerolog.SetGlobalLevel(logLevel)

	return zlog.Logger.WithContext(ctx), nil
}
