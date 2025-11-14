package validation

import (
	"github.com/go-playground/validator/v10"
)

func NewValidator() (*validator.Validate, error) {
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
