package graceful

import (
	"errors"
)

var (
	ErrForcedSignal   = errors.New("forceful stop: signal")
	ErrForcedTimeout  = errors.New("forceful stop: timeout")
	ErrGracefulSignal = errors.New("graceful stop: signal")
	ErrGracefulExit   = errors.New("graceful stop: exit")
)
