package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/traefik/paerser/cli"
	"github.com/traefik/paerser/file"
	"github.com/traefik/paerser/flag"
	"github.com/traefik/paerser/parser"
)

type File struct {
	envPrefix string
}

func (f *File) Load(args []string, cmd *cli.Command) (bool, error) {
	parsedArgs, err := flag.Parse(args, cmd.Configuration)
	if err != nil {
		_ = cmd.PrintHelp(os.Stdout)

		return true, err
	}

	configFile := parsedArgs[fmt.Sprintf("%s.config", parser.DefaultRootName)]

	if configFile == "" {
		if f.envPrefix != "" {
			configFile = os.Getenv(f.envPrefix + "CONFIG")
		}
	}

	if configFile == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return true, err
		}

		cfgPath := filepath.Join(cmd.Name, cmd.Name)

		finder := cli.Finder{
			BasePaths: []string{
				filepath.Join("/etc", cfgPath),
				filepath.Join(userConfigDir, cfgPath),
				cmd.Name,
			},
			Extensions: []string{"toml", "yaml", "yml"},
		}

		configFile, err = finder.Find("")
		if err != nil {
			return true, err
		}
	}

	if configFile == "" {
		return false, nil
	}

	if err := file.Decode(configFile, cmd.Configuration); err != nil {
		return true, fmt.Errorf("failed to decode configuration: %w ", err)
	}

	return false, nil
}
