package transaction

import (
	"testing"
)

func TestRollback(t *testing.T) {
	Rollback(t.Context())
}
