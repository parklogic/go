package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidHostname(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "letters", value: "example", want: true},
		{name: "digits and hyphen", value: "a1-b.example", want: true},
		{name: "empty", value: "", want: false},
		{name: "starts with hyphen", value: "-example", want: false},
		{name: "space", value: "ex ample", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, IsValidHostname(test.value))
		})
	}
}
