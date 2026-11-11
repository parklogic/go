package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubError struct{}

func (stubError) Error() string { return "boom" }

func TestNewError(t *testing.T) {
	t.Parallel()

	base := stubError{}

	tests := []struct {
		name    string
		err     error
		message string
	}{
		{
			name:    "wraps the cause",
			err:     base,
			message: "value error: boom",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require := require.New(t)
			assert := assert.New(t)

			err := NewError(test.err)
			require.NotNil(err)
			require.EqualError(err, test.message)
			require.ErrorIs(err, test.err)
			assert.Equal(test.err, err.Unwrap())
		})
	}
}
