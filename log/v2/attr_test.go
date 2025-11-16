package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAttr(t *testing.T) {
	tests := []struct {
		name  string
		attrs []any
		want  any
	}{
		{
			name: "No attributes",
			want: log,
		},

		{
			name:  "Cause with running context",
			attrs: []any{Cause(ctx)},
			want:  log,
		},
		{
			name:  "Cause with cancelled context",
			attrs: []any{Cause(ctxCancelled)},
			want: struct {
				baseLog
				Cause string `json:"cause"`
			}{
				baseLog: log,
				Cause:   context.Canceled.Error(),
			},
		},
		{
			name:  "Cause with cancelled context with cause",
			attrs: []any{Cause(ctxCancelledCause)},
			want: struct {
				baseLog
				Cause string `json:"cause"`
			}{
				baseLog: log,
				Cause:   cause.Error(),
			},
		},

		{
			name:  "Err with nil error",
			attrs: []any{Err(nil)},
			want:  log,
		},
		{
			name:  "Err with simple error",
			attrs: []any{Err(errSimple)},
			want: struct {
				baseLog
				Err string `json:"err"`
			}{
				baseLog: log,
				Err:     errSimple.Error(),
			},
		},

		{
			name:  "ErrWithType with nil error",
			attrs: []any{ErrWithType(nil)},
			want:  log,
		},
		{
			name:  "ErrWithType with simple error",
			attrs: []any{ErrWithType(errSimple)},
			want: struct {
				baseLog
				Err struct {
					Msg  string `json:"msg"`
					Type string `json:"type"`
				} `json:"err"`
			}{
				baseLog: log,
				Err: struct {
					Msg  string `json:"msg"`
					Type string `json:"type"`
				}{
					Msg:  errSimple.Error(),
					Type: "*errors.errorString",
				},
			},
		},
		{
			name:  "ErrWithType with custom error",
			attrs: []any{ErrWithType(errCustom)},
			want: struct {
				baseLog
				Err struct {
					Msg  string `json:"msg"`
					Type string `json:"type"`
				} `json:"err"`
			}{
				baseLog: log,
				Err: struct {
					Msg  string `json:"msg"`
					Type string `json:"type"`
				}{
					Msg:  errCustom.Error(),
					Type: "log.customError",
				},
			},
		},
	}

	for _, tt := range tests {
		b := &bytes.Buffer{}
		logger := slog.New(slog.NewJSONHandler(b, nil))
		want := marshalLog(t, tt.want)

		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				logger.Log(context.Background(), logLvl, logMsg, tt.attrs...)
				assert.Equal(t, want, b.String())
			})
		})
	}
}

func TestStackTraceSingleLevel(t *testing.T) {
	r := firstLevel()

	assert.Contains(t, r, "log/v2/attr.go:")
	assert.Contains(t, r, "github.com/parklogic/go/log/v2.stackTrace()")
	assert.Contains(t, r, "log/v2/attr_test.go:")
	assert.Contains(t, r, "github.com/parklogic/go/log/v2.firstLevel(...)")
	assert.NotContains(t, r, "secondLevel")
}

func TestStackTraceMultiLevel(t *testing.T) {
	r := secondLevel()

	assert.Contains(t, r, "log/v2/attr.go:")
	assert.Contains(t, r, "github.com/parklogic/go/log/v2.stackTrace()")
	assert.Contains(t, r, "log/v2/attr_test.go:")
	assert.Contains(t, r, "github.com/parklogic/go/log/v2.firstLevel(...)")
	assert.Contains(t, r, "github.com/parklogic/go/log/v2.secondLevel(...)")
}

type baseLog struct {
	Time  time.Time  `json:"time"`
	Level slog.Level `json:"level"`
	Msg   string     `json:"msg"`
}

var (
	logLvl = slog.LevelInfo
	logMsg = "testing"
	now    = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	log = baseLog{
		Time:  now.Local(),
		Level: logLvl,
		Msg:   logMsg,
	}
)

func marshalLog(t *testing.T, log any) string {
	b, err := json.Marshal(log)
	if err != nil {
		t.Fatal(err)
	}

	return string(b) + "\n"
}

type customError struct {
	msg string
}

func (e customError) Error() string {
	return e.msg
}

var (
	ctx                   = context.Background()
	ctxCancelled, c       = context.WithCancel(context.Background())
	cause                 = fmt.Errorf("cancel cause")
	ctxCancelledCause, cc = context.WithCancelCause(context.Background())
	errSimple             = fmt.Errorf("simple error")
	errCustom             = customError{msg: "custom error"}
)

func init() {
	c()
	cc(cause)
}

func firstLevel() string {
	return stackTrace()
}

func secondLevel() string {
	return firstLevel()
}
