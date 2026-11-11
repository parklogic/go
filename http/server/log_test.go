package server

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubAddr struct{}

func (stubAddr) Network() string { return "tcp" }

func (stubAddr) String() string { return "127.0.0.1:9" }

type stubListener struct{}

func (stubListener) Accept() (net.Conn, error) { return nil, io.EOF }

func (stubListener) Close() error { return nil }

func (stubListener) Addr() net.Addr { return stubAddr{} }

func TestHTTPBaseContext(t *testing.T) {
	t.Parallel()

	base := HTTPBaseContext(t.Context())
	got := base(stubListener{})
	require.NotNil(t, got)
}

func TestHTTPConnContext(t *testing.T) {
	t.Parallel()

	left, right := net.Pipe()

	t.Cleanup(func() {
		require.NoError(t, left.Close())
		require.NoError(t, right.Close())
	})

	got := HTTPConnContext(t.Context(), left)
	require.NotNil(t, got)
}

func TestHTTPConnState(t *testing.T) {
	assert.Nil(t, HTTPConnState(t.Context()))

	var written bytes.Buffer

	logger := zerolog.New(&written).Level(zerolog.TraceLevel)
	observe := HTTPConnState(logger.WithContext(t.Context()))
	require.NotNil(t, observe)

	left, right := net.Pipe()

	t.Cleanup(func() {
		require.NoError(t, left.Close())
		require.NoError(t, right.Close())
	})

	observe(left, http.StateNew)
	observe(left, http.StateClosed)
	observe(left, http.StateIdle)
	assert.Contains(t, written.String(), "New connection")
	assert.Contains(t, written.String(), "Connection closed")
}

func TestHTTPLogger(t *testing.T) {
	t.Parallel()

	require.NotNil(t, HTTPLogger(t.Context()))
}
