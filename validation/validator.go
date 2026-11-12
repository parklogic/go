package validation

import (
	"github.com/go-playground/validator/v10"
)

// registered is the validator New stores and Ctx returns on a miss.
//
//nolint:gochecknoglobals // one value; a failed registration panics while the package loads.
var registered = mustValidator()

func mustValidator() *validator.Validate {
	v, err := newValidator()
	if err != nil {
		panic(err)
	}

	return v
}

func NewValidator() (*validator.Validate, error) {
	return registered, nil
}

func newValidator() (*validator.Validate, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	for tag, fn := range map[string]validator.FuncCtx{
		"fqdn_allow_wildcard": fqdnAllowWildcard,
	} {
		if err := v.RegisterValidationCtx(tag, fn); err != nil {
			return nil, err
		}
	}

	return v, nil
}
