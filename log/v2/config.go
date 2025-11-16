package log

import (
	"fmt"
)

// Configuration is used by [github.com/traefik/paerser] to load the logger configuration.
type Configuration struct {
	AddSource bool   `description:"Add caller source location to log output"`
	Format    string `description:"Log format output (\"logfmt\", \"json\")"`
	Level     string `description:"Logging level (\"debug\", \"info\", \"warn\", \"error\")"`
}

// NewConfiguration returns the default configuration for the logger.
func NewConfiguration() *Configuration {
	return &Configuration{
		AddSource: true,
		Format:    "json",
		Level:     "info",
	}
}

type ErrInvalidConfig struct {
	Field string
	Value string
}

func (e ErrInvalidConfig) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Value)
}
