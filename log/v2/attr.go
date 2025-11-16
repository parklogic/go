package log

import (
	"log/slog"
	"reflect"
	"runtime"
)

// Err returns a [slog.Attr] for an error message.
func Err(value error) slog.Attr {
	return slog.String("err", value.Error())
}

// ErrWithType returns a [slog.Attr] for an error message and its type.
func ErrWithType(value error) slog.Attr {
	return slog.Group("err",
		slog.String("msg", value.Error()),
		slog.String("type", reflect.TypeOf(value).String()),
	)
}

// ErrWithStack returns a [slog.Attr] for an error message, its type and a stack trace of the call site.
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
