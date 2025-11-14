package publicsuffix

type Configuration struct {
	CacheSize    int    `description:"Size of the parsed domains cache"`
	ErrCacheSize int    `description:"Size of the error cache"`
	Path         string `description:"Path to the public suffix list file"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		CacheSize:    1024,
		ErrCacheSize: 128,
		Path:         DefaultListPath,
	}
}
