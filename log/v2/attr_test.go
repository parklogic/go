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

type baseLog struct {
	Time  time.Time  `json:"time"`
	Level slog.Level `json:"level"`
	Msg   string     `json:"msg"`
}

var synctestNow = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

func marshalLog(t *testing.T, log any) string {
	b, err := json.Marshal(log)
	if err != nil {
		t.Fatal(err)
	}

	return string(b) + "\n"
}

func TestCause(t *testing.T) {
	logLvl := slog.LevelInfo
	logMsg := "testing Cause"

	ctx := context.Background()

	ctxCancelled, c := context.WithCancel(context.Background())
	c()

	cause := fmt.Errorf("cancel cause")

	ctxCancelledCause, cc := context.WithCancelCause(context.Background())
	cc(cause)

	tests := []struct {
		name string
		ctx  context.Context
		want any
	}{
		{
			name: "Running context",
			ctx:  ctx,
			want: baseLog{
				Time:  synctestNow.Local(),
				Level: logLvl,
				Msg:   logMsg,
			},
		},
		{
			name: "Without cause",
			ctx:  ctxCancelled,
			want: struct {
				baseLog
				Cause string `json:"cause"`
			}{
				baseLog: baseLog{
					Time:  synctestNow.Local(),
					Level: logLvl,
					Msg:   logMsg,
				},
				Cause: context.Canceled.Error(),
			},
		},
		{
			name: "With cause",
			ctx:  ctxCancelledCause,
			want: struct {
				baseLog
				Cause string `json:"cause"`
			}{
				baseLog: baseLog{
					Time:  synctestNow.Local(),
					Level: logLvl,
					Msg:   logMsg,
				},
				Cause: cause.Error(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				b := &bytes.Buffer{}
				logger := slog.New(slog.NewJSONHandler(b, nil))
				logger.Log(test.ctx, logLvl, logMsg, Cause(test.ctx))

				want := marshalLog(t, test.want)

				assert.Equal(t, want, b.String())
			})
		})
	}
}

var errSimple = fmt.Errorf("simple error")

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

var errCustom = &customError{msg: "custom error"}

func TestErr(t *testing.T) {
	logLvl := slog.LevelInfo
	logMsg := "testing Err"

	tests := []struct {
		name string
		err  error
		want any
	}{
		{
			name: "Simple error",
			err:  errSimple,
			want: struct {
				baseLog
				Err string `json:"err"`
			}{
				baseLog: baseLog{
					Time:  synctestNow.Local(),
					Level: logLvl,
					Msg:   logMsg,
				},
				Err: errSimple.Error(),
			},
		},
		{
			name: "nil error",
			err:  nil,
			want: baseLog{
				Time:  synctestNow.Local(),
				Level: logLvl,
				Msg:   logMsg,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				b := &bytes.Buffer{}
				logger := slog.New(slog.NewJSONHandler(b, nil))
				logger.Log(context.Background(), logLvl, logMsg, Err(tt.err))

				want := marshalLog(t, tt.want)

				assert.Equal(t, want, b.String())
			})
		})
	}
}

func TestErrWithType(t *testing.T) {
	logLvl := slog.LevelInfo
	logMsg := "testing ErrWithType"

	tests := []struct {
		name string
		err  error
		want any
	}{
		{
			name: "simple error",
			err:  errSimple,
			want: struct {
				baseLog
				Err struct {
					Msg  string `json:"msg"`
					Type string `json:"type"`
				} `json:"err"`
			}{
				baseLog: baseLog{
					Time:  synctestNow.Local(),
					Level: logLvl,
					Msg:   logMsg,
				},
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
			name: "custom error",
			err:  errCustom,
			want: struct {
				baseLog
				Err struct {
					Msg  string `json:"msg"`
					Type string `json:"type"`
				} `json:"err"`
			}{
				baseLog: baseLog{
					Time:  synctestNow.Local(),
					Level: logLvl,
					Msg:   logMsg,
				},
				Err: struct {
					Msg  string `json:"msg"`
					Type string `json:"type"`
				}{
					Msg:  errCustom.Error(),
					Type: "*log.customError",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				b := &bytes.Buffer{}
				logger := slog.New(slog.NewJSONHandler(b, nil))
				logger.Log(context.Background(), logLvl, logMsg, ErrWithType(tt.err))

				want := marshalLog(t, tt.want)

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

func firstLevel() string {
	return stackTrace()
}

func secondLevel() string {
	return firstLevel()
}
