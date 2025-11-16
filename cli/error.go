package cli

type ErrInvalidConfig interface {
	error
	Field() string
	Value() string
}
