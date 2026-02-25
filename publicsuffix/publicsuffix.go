package publicsuffix

import (
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
	"golang.org/x/net/idna"
)

var DefaultFindOptions = &publicsuffix.FindOptions{
	IgnorePrivate: true,
	DefaultRule:   publicsuffix.DefaultRule,
}

var NoStarFindOptions = &publicsuffix.FindOptions{
	IgnorePrivate: true,
	DefaultRule:   nil,
}

const DefaultListPath = "/usr/share/publicsuffix/public_suffix_list.dat"

var DefaultParserOptions = &publicsuffix.ParserOption{
	PrivateDomains: false,
	ASCIIEncoded:   false,
}

func Expand(d *publicsuffix.DomainName) (fqdn, domain, tld, subdomain string) {
	tld = d.TLD

	domain = d.TLD
	if d.SLD != "" {
		domain = d.SLD + "." + domain
	}

	subdomain = d.TRD

	fqdn = d.String()

	return fqdn, domain, tld, subdomain
}

func Normalise(d *publicsuffix.DomainName) error {
	tld, err := idna.Lookup.ToASCII(d.TLD)
	if err != nil {
		return err
	}

	d.TLD = tld

	sld, err := idna.Lookup.ToASCII(d.SLD)
	if err != nil {
		return err
	}

	d.SLD = sld

	var hasWildcard bool

	if d.TRD == "*" || strings.HasPrefix(d.TRD, "*.") {
		hasWildcard = true
		d.TRD = strings.TrimPrefix(strings.TrimPrefix(d.TRD, "*"), ".")
	}

	trd, err := idna.Lookup.ToASCII(d.TRD)
	if err != nil {
		return err
	}

	if hasWildcard {
		if trd == "" {
			trd = "*"
		} else {
			trd = "*" + "." + trd
		}
	}

	d.TRD = trd

	return nil
}
