package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type hostName struct {
	Name string `validate:"fqdn_allow_wildcard"`
}

func TestNewValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "fqdn", value: "example.com"},
		{name: "wildcard", value: "*.example.com"},
		{name: "empty", value: "", wantErr: true},
		{name: "leading hyphen", value: "-bad.com", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require := require.New(t)

			validator, err := NewValidator()
			require.NoError(err)

			err = validator.StructCtx(t.Context(), hostName{Name: test.value})
			if test.wantErr {
				require.Error(err)

				return
			}

			require.NoError(err)
		})
	}
}
