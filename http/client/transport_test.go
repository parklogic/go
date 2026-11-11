package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const preservedAgent = "keep-me"

func TestNewTransport(t *testing.T) {
	t.Parallel()

	cfg := NewConfiguration()
	got := newTransport(cfg)
	require.Equal(t, cfg.UserAgent, got.userAgent)
	assert.Equal(t, cfg.IdleConnTimeout, got.transport.IdleConnTimeout)
	assert.Equal(t, cfg.ResponseHeaderTimeout, got.transport.ResponseHeaderTimeout)
	assert.Equal(t, cfg.TLSHandshakeTimeout, got.transport.TLSHandshakeTimeout)
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)

	cfg := NewConfiguration()
	transport := newTransport(cfg)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
	res, err := transport.RoundTrip(req)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, res.Body.Close()) })
	require.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, cfg.UserAgent, req.Header.Get("User-Agent"))

	kept := httptest.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
	kept.Header.Set("User-Agent", preservedAgent)
	keptRes, err := transport.RoundTrip(kept)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, keptRes.Body.Close()) })
	assert.Equal(t, preservedAgent, kept.Header.Get("User-Agent"))

	canceled, cancel := context.WithCancel(t.Context())
	cancel()

	failed := httptest.NewRequestWithContext(canceled, http.MethodGet, srv.URL, nil)
	failedRes, err := transport.RoundTrip(failed)
	require.Error(t, err)

	if failedRes != nil {
		require.NoError(t, failedRes.Body.Close())
	}
}
