package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerMethodNotAllowed(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	_, _, err := HandlerMethodNotAllowed(httptest.NewRecorder(), req, nil)
	require.ErrorIs(t, err, ErrMethodNotAllowed)

	var apiErr Error

	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusMethodNotAllowed, apiErr.StatusCode())
}

func TestHandlerNotFound(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	_, _, err := HandlerNotFound(httptest.NewRecorder(), req, nil)
	require.ErrorIs(t, err, ErrEndpointNotFound)
}

func TestHandlerLive(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/livez", nil)
	status, body, err := HandlerLive(httptest.NewRecorder(), req, NilPayload{})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "ok", body.Status)
}
