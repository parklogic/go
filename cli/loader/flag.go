package loader

import (
	"fmt"
	"strings"

	"github.com/traefik/paerser/cli"
	"github.com/traefik/paerser/flag"
)

type Flag struct{}

func (*Flag) Load(args []string, cmd *cli.Command) (bool, error) {
	args = withoutConfigArg(args)
	if len(args) == 0 {
		return false, nil
	}

	if err := flag.Decode(args, cmd.Configuration); err != nil {
		return true, fmt.Errorf("failed to decode configuration from flags: %w", err)
	}

	return false, nil
}

// withoutConfigArg drops the file loader's --config flag.
// flag.Decode rejects it because it is not a configuration field.
func withoutConfigArg(args []string) []string {
	kept := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--config" {
			i++

			continue
		}

		if strings.HasPrefix(arg, "--config=") {
			continue
		}

		kept = append(kept, arg)
	}

	return kept
}
