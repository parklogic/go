package log

import (
	"context"
	"log/slog"
	"reflect"
	"runtime"
)

// Cause returns a [slog.Attr] for the cause of the context cancellation, or an empty attribute if the context was not cancelled.
func Cause(ctx context.Context) slog.Attr {
	cause := context.Cause(ctx)
	if cause == nil {
		return slog.Attr{}
	}

	return slog.String("cause", cause.Error())
}

// Err returns a [slog.Attr] for an error message, or an empty attribute if value is nil.
func Err(value error) slog.Attr {
	if value == nil {
		return slog.Attr{}
	}

	return slog.String("err", value.Error())
}

// ErrWithType returns a [slog.Attr] for an error message and its type, or an empty attribute if value is nil.
func ErrWithType(value error) slog.Attr {
	if value == nil {
		return slog.Attr{}
	}

	return slog.Group("err",
		slog.String("msg", value.Error()),
		slog.String("type", reflect.TypeOf(value).String()),
	)
}

// ErrWithStack returns a [slog.Attr] for an error message, its type, and a stack trace of the call site, or an empty attribute if value is nil.
func ErrWithStack(value error) slog.Attr {
	return slog.Group("err",
		slog.String("msg", value.Error()),
		slog.String("type", reflect.TypeOf(value).String()),
		slog.String("stack", stackTrace()),
	)
}

// stackTrace returns a stack trace of the call site.
func stackTrace() string {
	buf := make([]byte, 1024)

	s := runtime.Stack(buf, false)
	if s == 0 {
		return ""
	}

	return string(buf[:s])
}
