package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtx(t *testing.T) {
	t.Parallel()

	assert.Equal(t, http.DefaultClient, Ctx(t.Context()))
}

func TestWithContext(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	require.Equal(t, client, Ctx(WithContext(t.Context(), client)))
	assert.Nil(t, Ctx(WithContext(t.Context(), nil)))
}
