package cli

type Configuration struct {
	Config string `description:"Path to the configuration file (tries to load from multiple locations if not set)"`
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}
