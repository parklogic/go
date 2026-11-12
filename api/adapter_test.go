package api

import (
	"context"
	"encoding/json"
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

type nullBody struct{}

func (nullBody) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
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

func TestAdapterEmptyBody(t *testing.T) {
	t.Parallel()

	handler := Adapter(func(http.ResponseWriter, *http.Request, namedPayload) (int, *namedPayload, error) {
		return http.StatusNoContent, nil, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Empty(t, rec.Body.String())
}

func TestAdapterRejectsUnexpectedPayload(t *testing.T) {
	t.Parallel()

	handler := Adapter(func(http.ResponseWriter, *http.Request, NilPayload) (int, *struct{}, error) {
		return http.StatusOK, nil, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader("{}"))
	req.ContentLength = 2
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "No payload expected")
}

func TestAdapterProblemJSON(t *testing.T) {
	t.Parallel()

	handler := Adapter(func(http.ResponseWriter, *http.Request, struct{}) (int, *struct{}, error) {
		return 0, nil, stubAPIError{}
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "Internal server error")
}

func TestAdapterCanceled(t *testing.T) {
	t.Parallel()

	handler := Adapter(func(http.ResponseWriter, *http.Request, struct{}) (int, *struct{}, error) {
		return 0, nil, context.Canceled
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Empty(t, rec.Body.String())
	require.Empty(t, rec.Header().Get("Content-Type"))
}

func TestAdapterNullBody(t *testing.T) {
	t.Parallel()

	handler := Adapter(func(http.ResponseWriter, *http.Request, struct{}) (int, *nullBody, error) {
		return http.StatusOK, &nullBody{}, nil
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body.String())
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestAdapterAbortHandler(t *testing.T) {
	t.Parallel()

	handler := Adapter(func(http.ResponseWriter, *http.Request, struct{}) (int, *struct{}, error) {
		panic(http.ErrAbortHandler)
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	require.PanicsWithValue(t, http.ErrAbortHandler, func() {
		handler.ServeHTTP(httptest.NewRecorder(), req)
	})
}

func TestAdapterJSONRoundTrip(t *testing.T) {
	t.Parallel()

	var got namedPayload

	handler := Adapter(func(_ http.ResponseWriter, _ *http.Request, payload namedPayload) (int, *namedPayload, error) {
		got = payload

		return http.StatusOK, &namedPayload{Name: payload.Name}, nil
	})

	const body = `{"name":"a"}`

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(body))
	req.ContentLength = int64(len(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "a", got.Name)

	var decoded namedPayload

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &decoded))
	require.Equal(t, got, decoded)
}
