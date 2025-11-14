package log

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

type TracedReader struct {
	r io.Reader
	l *zerolog.Logger
}

func NewTracedReader(ctx context.Context, r io.Reader) io.Reader {
	logger := zerolog.Ctx(ctx)

	if logger.GetLevel() > zerolog.TraceLevel {
		return r
	}

	return &TracedReader{
		r: r,
		l: logger,
	}
}

func (r *TracedReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)

	if err == nil {
		r.l.Trace().Bytes("data", p[:n]).Msg("Data received")
	}

	return n, err
}

type TracedWriter struct {
	w io.Writer
	l *zerolog.Logger
}

func NewTracedWriter(ctx context.Context, w io.Writer) io.Writer {
	logger := zerolog.Ctx(ctx)

	if logger.GetLevel() > zerolog.TraceLevel {
		return w
	}

	return &TracedWriter{
		w: w,
		l: logger,
	}
}

func (w *TracedWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)

	if err == nil {
		w.l.Trace().Bytes("data", p).Msg("Data sent")
	}

	return n, err
}
