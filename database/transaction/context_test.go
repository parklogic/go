package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCtx(t *testing.T) {
	assert.Nil(t, Ctx(t.Context()))
}

func TestWithContext(t *testing.T) {
	ctx := WithContext(t.Context(), nil)
	assert.Nil(t, Ctx(ctx))
}
