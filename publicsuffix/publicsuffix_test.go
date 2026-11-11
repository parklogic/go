package publicsuffix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"golang.org/x/net/idna"
)

const suffixTLD = "com"

func TestExpand(t *testing.T) {
	t.Parallel()

	name := &publicsuffix.DomainName{TLD: suffixTLD, SLD: "example", TRD: "www"}
	fqdn, domain, tld, subdomain := Expand(name)
	assert.Equal(t, name.String(), fqdn)
	assert.Equal(t, "example."+suffixTLD, domain)
	assert.Equal(t, suffixTLD, tld)
	assert.Equal(t, "www", subdomain)
}

func TestNormalise(t *testing.T) {
	t.Parallel()

	name := &publicsuffix.DomainName{TLD: suffixTLD, SLD: "exämple", TRD: "*.www"}
	require.NoError(t, Normalise(name))

	wantSLD, err := idna.Lookup.ToASCII("exämple")
	require.NoError(t, err)
	assert.Equal(t, wantSLD, name.SLD)
	assert.Equal(t, "*.www", name.TRD)

	invalid := &publicsuffix.DomainName{TLD: suffixTLD, SLD: "bad name"}
	require.Error(t, Normalise(invalid))
}
