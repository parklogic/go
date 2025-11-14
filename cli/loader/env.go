package loader

import (
	"fmt"
	"os"

	"github.com/traefik/paerser/cli"
	"github.com/traefik/paerser/env"
)

type Env struct {
	prefix string
}

func (e *Env) Load(_ []string, cmd *cli.Command) (bool, error) {

	vars := env.FindPrefixedEnvVars(os.Environ(), e.prefix, cmd.Configuration)
	if len(vars) == 0 {
		return false, nil
	}

	if err := env.Decode(vars, e.prefix, cmd.Configuration); err != nil {
		return true, fmt.Errorf("failed to decode configuration from environment variables: %w ", err)
	}

	return false, nil
}
