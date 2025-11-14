package loader

import (
	"regexp"
	"strings"

	"github.com/traefik/paerser/cli"
)

var invalidEnvPrefixChars = regexp.MustCompile(`[^a-zA-Z0-9]`)

func New(envPrefix string) []cli.ResourceLoader {
	envPrefix = cleanEnvPrefix(envPrefix)

	return []cli.ResourceLoader{
		&File{envPrefix: envPrefix},
		&Env{prefix: envPrefix},
		&Flag{},
	}
}

func cleanEnvPrefix(prefix string) string {
	return strings.ToUpper(invalidEnvPrefixChars.ReplaceAllString(prefix, "")) + "_"
}
