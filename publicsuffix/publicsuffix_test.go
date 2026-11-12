package publicsuffix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/idna"
)

const suffixTLD = "com"

func TestParseNormalisesLabels(t *testing.T) {
	t.Parallel()

	list, err := NewList(writeList(t), 0, 0)
	require.NoError(t, err)

	got, err := list.Parse("www.exämple.com")
	require.NoError(t, err)

	wantSLD, err := idna.Lookup.ToASCII("exämple")
	require.NoError(t, err)
	assert.Equal(t, wantSLD+".com", got.Domain)
	assert.Equal(t, "www", got.Subdomain)
	assert.Equal(t, suffixTLD, got.TLD)

	_, err = list.Parse("bad name.com")
	require.Error(t, err)
}
