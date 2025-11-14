package value

import (
	"fmt"
)

type Error struct {
	err error
}

func NewError(err error) *Error {
	return &Error{
		err: err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("value error: %s", e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}
