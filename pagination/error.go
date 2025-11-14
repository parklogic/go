package pagination

import (
	"fmt"
)

var (
	ErrInvalidLimit     = fmt.Errorf("invalid limit")
	ErrInvalidPage      = fmt.Errorf("invalid page number")
	ErrInvalidSortKey   = fmt.Errorf("invalid sort key")
	ErrInvalidSortOrder = fmt.Errorf("invalid sort order")
	ErrMissingSortKey   = fmt.Errorf("missing sort key")
)
