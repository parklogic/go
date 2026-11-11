package bearer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type staticValidator struct {
	valid bool
	err   error
}

func (v staticValidator) Validate(context.Context, string) (bool, error) {
	return v.valid, v.err
}

type stubBearerError struct{}

func (stubBearerError) Error() string { return "validator failed" }

const bearerHeader = "Bearer secret"

func okNext(_ http.ResponseWriter, _ *http.Request, _ struct{}) (int, *struct{}, error) {
	status := http.StatusNoContent

	return status, &struct{}{}, nil
}

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		validator Validator
		wantErr   error
	}{
		{name: "missing header", wantErr: ErrMissingHeader},
		{name: "wrong scheme", header: "Basic secret", wantErr: ErrInvalidAuthType},
		{name: "missing token", header: "Bearer", wantErr: ErrMissingToken},
		{name: "rejected token", header: bearerHeader, validator: staticValidator{}, wantErr: ErrInvalidToken},
		{
			name:      "validator error",
			header:    bearerHeader,
			validator: staticValidator{err: stubBearerError{}},
			wantErr:   stubBearerError{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
			if test.header != "" {
				req.Header.Set("Authorization", test.header)
			}

			if test.validator != nil {
				req = req.WithContext(WithContext(req.Context(), test.validator))
			}

			_, body, err := Middleware(okNext)(httptest.NewRecorder(), req, struct{}{})
			require.ErrorIs(err, test.wantErr)
			assert.Nil(body)
		})
	}
}

func TestMiddlewareAcceptsToken(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("Authorization", bearerHeader)
	req = req.WithContext(WithContext(req.Context(), staticValidator{valid: true}))

	status, body, err := Middleware(okNext)(httptest.NewRecorder(), req, struct{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, status)
	assert.NotNil(t, body)
}

func TestMiddlewareMissingValidator(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	req.Header.Set("Authorization", bearerHeader)

	_, body, err := Middleware(okNext)(httptest.NewRecorder(), req, struct{}{})
	require.Error(t, err)
	assert.Nil(t, body)
}
