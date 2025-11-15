package linode

import (
	"context"

	"github.com/linode/go-metadata"
	"github.com/rs/zerolog"
)

type clientLogger struct {
	l *zerolog.Logger
}

func (c clientLogger) Errorf(format string, v ...interface{}) {
	c.l.Error().Msgf(format, v...)
}

func (c clientLogger) Warnf(format string, v ...interface{}) {
	c.l.Warn().Msgf(format, v...)
}

func (c clientLogger) Debugf(format string, v ...interface{}) {
	c.l.Trace().Msgf(format, v...)
}

func NewClient(ctx context.Context) (*metadata.Client, error) {
	return metadata.NewClient(ctx,
		metadata.ClientWithLogger(clientLogger{zerolog.Ctx(ctx)}),
		metadata.ClientWithDebug(),
	)
}
