package pprof

type Configuration struct {
	Address string `description:"Address for the profiler to listen on"`
	Enabled bool   `description:"Enable profiling"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Address: ":6060",
	}
}
