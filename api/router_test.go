package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mychi "github.com/parklogic/go/chi"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	router := NewRouter(mychi.NewConfiguration())
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/livez", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body Status

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "ok", body.Status)
}
