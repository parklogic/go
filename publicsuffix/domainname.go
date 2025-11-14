package publicsuffix

import (
	"strings"
)

type DomainName struct {
	FQDN      string
	Domain    string
	Subdomain string
	TLD       string
}

func (d *DomainName) String() string {
	return d.FQDN
}

func (d *DomainName) Wildcard() string {
	if d.Subdomain == "" {
		return d.Domain
	}

	_, parentSubdomain, hasParentSubdomain := strings.Cut(d.Subdomain, ".")

	if !hasParentSubdomain {
		return "*." + d.Domain
	}

	return "*." + parentSubdomain + "." + d.Domain
}
