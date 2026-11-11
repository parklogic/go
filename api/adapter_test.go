package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type namedPayload struct {
	Name string `json:"name"`
}

func TestAdapter(t *testing.T) {
	t.Parallel()

	handler := Adapter(func(_ http.ResponseWriter, _ *http.Request, payload namedPayload) (int, *namedPayload, error) {
		return http.StatusCreated, &namedPayload{Name: payload.Name}, nil
	})

	const body = `{"name":"a"}`

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(body))
	req.ContentLength = int64(len(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, body, rec.Body.String())
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestDecodePayload(t *testing.T) {
	t.Parallel()

	empty := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	require.NoError(t, decodePayload(empty, &namedPayload{}))

	rejected := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader("{}"))
	rejected.ContentLength = 2
	require.ErrorIs(t, decodePayload(rejected, &NilPayload{}), ErrPayloadNotExpected)

	const body = `{"name":"a"}`

	accepted := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(body))
	accepted.ContentLength = int64(len(body))

	var got namedPayload

	require.NoError(t, decodePayload(accepted, &got))
	require.Equal(t, "a", got.Name)
}

func TestExecHandlerFunc(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	handler := func(http.ResponseWriter, *http.Request, struct{}) (int, *struct{}, error) {
		return http.StatusNoContent, nil, nil
	}
	status, body, err := execHandlerFunc(httptest.NewRecorder(), req, handler)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)
	require.Nil(t, body)

	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		abort := func(http.ResponseWriter, *http.Request, struct{}) (int, *struct{}, error) {
			panic(http.ErrAbortHandler)
		}
		_, _, _ = execHandlerFunc(httptest.NewRecorder(), req, abort)
	})
}

func TestHandlerErr(t *testing.T) {
	t.Parallel()

	canceled := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handlerErr(rec, canceled, context.Canceled)
	require.Empty(t, rec.Body.String())

	failed := httptest.NewRecorder()
	handlerErr(failed, canceled, stubAPIError{})
	require.Equal(t, http.StatusInternalServerError, failed.Code)
	assert.Contains(t, failed.Body.String(), "Internal server error")
}
