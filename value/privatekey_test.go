package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPrivateKeyType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "rsa size", value: "2048", want: true},
		{name: "p256", value: "P256", want: true},
		{name: "unknown", value: "P384", want: false},
		{name: "empty", value: "", want: false},
		{name: "lowercase p256", value: "p256", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, IsValidPrivateKeyType(test.value))
		})
	}
}
