package publicsuffix

import (
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
	"golang.org/x/net/idna"
)

const DefaultListPath = "/usr/share/publicsuffix/public_suffix_list.dat"

var DefaultParserOptions = &publicsuffix.ParserOption{
	PrivateDomains: false,
	ASCIIEncoded:   false,
}

func findOptions(noStar bool) *publicsuffix.FindOptions {
	options := &publicsuffix.FindOptions{
		IgnorePrivate: true,
		DefaultRule:   publicsuffix.DefaultRule,
	}
	if noStar {
		options.DefaultRule = nil
	}

	return options
}

func expand(parsed *publicsuffix.DomainName) DomainName {
	tld := parsed.TLD
	domain := tld

	if parsed.SLD != "" {
		domain = parsed.SLD + "." + domain
	}

	return DomainName{
		FQDN:      parsed.String(),
		Domain:    domain,
		Subdomain: parsed.TRD,
		TLD:       tld,
	}
}

func normalise(parsed *publicsuffix.DomainName) error {
	tld, err := idna.Lookup.ToASCII(parsed.TLD)
	if err != nil {
		return err
	}

	parsed.TLD = tld

	sld, err := idna.Lookup.ToASCII(parsed.SLD)
	if err != nil {
		return err
	}

	parsed.SLD = sld

	var hasWildcard bool

	if parsed.TRD == "*" || strings.HasPrefix(parsed.TRD, "*.") {
		hasWildcard = true
		parsed.TRD = strings.TrimPrefix(strings.TrimPrefix(parsed.TRD, "*"), ".")
	}

	trd, err := idna.Lookup.ToASCII(parsed.TRD)
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

	parsed.TRD = trd

	return nil
}
