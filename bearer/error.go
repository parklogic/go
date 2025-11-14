package bearer

import (
	"fmt"
)

var (
	ErrInvalidAuthType = fmt.Errorf("invalid authorization type")
	ErrInvalidToken    = fmt.Errorf("invalid bearer token")
	ErrMissingHeader   = fmt.Errorf("missing authorization header")
	ErrMissingToken    = fmt.Errorf("missing bearer token")
)
