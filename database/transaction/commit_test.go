package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommit(t *testing.T) {
	err := Commit(t.Context())
	require.NoError(t, err)
}
