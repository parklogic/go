package pagination

type Configuration struct {
	LimitParameter string
	PageParameter  string
	SortParameter  string

	DefaultSort string

	DefaultLimit int64
	MaxLimit     int64
	MinLimit     int64
}

func newConfiguration(f ...func(*Configuration)) Configuration {
	config := Configuration{
		LimitParameter: "limit",
		PageParameter:  "page",
		SortParameter:  "sort",

		DefaultSort: "id",

		DefaultLimit: 25,
		MaxLimit:     1000,
		MinLimit:     0,
	}

	for _, ff := range f {
		ff(&config)
	}

	return config
}

func WithDefaultSize(size int64) func(*Configuration) {
	return func(config *Configuration) {
		config.DefaultLimit = size
	}
}

func WithDefaultSort(sort string) func(*Configuration) {
	return func(config *Configuration) {
		config.DefaultSort = sort
	}
}

func WithMaxSize(size int64) func(*Configuration) {
	return func(config *Configuration) {
		config.MaxLimit = size
	}
}

func WithMinSize(size int64) func(*Configuration) {
	return func(config *Configuration) {
		config.MinLimit = size
	}
}

func WithOptionalPagination() func(*Configuration) {
	return WithMinSize(-1)
}
