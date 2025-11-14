package log

type Configuration struct {
	ForceColor   bool   `description:"Force color output (defaults to JSON output if stderr is not a TTY)"`
	Level        string `description:"Logging level"`
	ShortenPath  bool   `description:"Shorten the caller file path in log output"`
	RemovePrefix string `description:"Remove prefix from caller file"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Level: "info",
	}
}
