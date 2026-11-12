# ParkLogic Go

Shared packages for ParkLogic. These names are the ones the public-suffix code should use.

## Language

**Domain name**:
A hostname parsed into its FQDN, domain, subdomain, and TLD.
_Avoid_: URL, host string

**No-star parse**:
A domain-name parse with the public suffix list's default `*` rule turned off.
_Avoid_: options parse, custom find options
