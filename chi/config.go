package chi

import (
	"time"
)

type Configuration struct {
	CompressionLevel int           `description:"Compression level for response bodies"`
	SlowResponse     time.Duration `description:"Slow response threshold"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		CompressionLevel: 5,
		SlowResponse:     200 * time.Millisecond,
	}
}
