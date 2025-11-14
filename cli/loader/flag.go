package loader

import (
	"fmt"

	"github.com/traefik/paerser/cli"
	"github.com/traefik/paerser/flag"
)

type Flag struct{}

func (*Flag) Load(args []string, cmd *cli.Command) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	if err := flag.Decode(args, cmd.Configuration); err != nil {
		return true, fmt.Errorf("failed to decode configuration from flags: %w", err)
	}

	return false, nil
}
