package log

import (
	"log"

	"github.com/rs/zerolog"
)

func setStdLogger(l *zerolog.Logger) {
	log.SetFlags(0)
	log.SetOutput(l)
}
