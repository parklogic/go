package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	var seen *Pagination

	next := func(_ http.ResponseWriter, r *http.Request, _ struct{}) (int, *struct{}, error) {
		seen = Ctx(r.Context())

		return http.StatusOK, &struct{}{}, nil
	}

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/items?limit=5&page=2", nil)
	status, _, err := Middleware(next)(httptest.NewRecorder(), req, struct{}{})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.NotNil(t, seen)
	require.Equal(t, int64(5), seen.Limit())
	require.Equal(t, int64(2), seen.Page())
}
