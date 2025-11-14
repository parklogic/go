package value

import (
	"regexp"
)

var (
	// validHostnameCharsRegexp is a regular expression to validate hostnames
	validHostnameCharsRegexp = regexp.MustCompile(`^([a-zA-Z0-9]{1}[a-zA-Z0-9-_.]{0,254})$`)
)

func IsValidHostname(v string) bool {
	return validHostnameCharsRegexp.MatchString(v)
}
