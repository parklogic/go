package chi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubPanicError struct{}

func (stubPanicError) Error() string { return "panic" }

func TestLogMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     zerolog.Level
		status    int
		threshold time.Duration
		want      string
	}{
		{name: "disabled", level: zerolog.Disabled, status: http.StatusOK, threshold: time.Second},
		{
			name: "failed", level: zerolog.ErrorLevel, status: http.StatusInternalServerError,
			threshold: time.Second, want: "Request failed",
		},
		{name: "slow", level: zerolog.WarnLevel, status: http.StatusOK, threshold: -time.Nanosecond, want: "Slow request"},
		{name: "missing", level: zerolog.TraceLevel, status: http.StatusNotFound, threshold: time.Hour, want: "Bad request"},
		{
			name: "rejected", level: zerolog.InfoLevel, status: http.StatusBadRequest,
			threshold: time.Hour, want: "Bad request",
		},
		{
			name: "completed", level: zerolog.TraceLevel, status: http.StatusOK,
			threshold: time.Hour, want: "Request completed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var written bytes.Buffer

			logger := zerolog.New(&written).Level(test.level)
			req := httptest.NewRequestWithContext(logger.WithContext(t.Context()), http.MethodGet, "/items", nil)
			rec := httptest.NewRecorder()
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(test.status)
			})

			LogMiddleware(test.threshold)(next).ServeHTTP(rec, req)
			require.Equal(t, test.status, rec.Code)

			if test.want == "" {
				assert.Empty(t, written.String())

				return
			}

			assert.Contains(t, written.String(), test.want)
		})
	}
}

func TestRecoverMiddleware(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(stubPanicError{})
	})
	RecoverMiddleware(next).ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)

	upgraded := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	upgraded.Header.Set("Connection", "Upgrade")

	skipped := httptest.NewRecorder()
	RecoverMiddleware(next).ServeHTTP(skipped, upgraded)
	require.Equal(t, http.StatusOK, skipped.Code)

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		abort := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			panic(http.ErrAbortHandler)
		})
		RecoverMiddleware(abort).ServeHTTP(httptest.NewRecorder(), req)
	})
}
