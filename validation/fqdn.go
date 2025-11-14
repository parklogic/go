package validation

import (
	"context"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var fqdnAllowWildcardRegex = regexp.MustCompile(`^(\*\.)?([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62}){1}(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?$`)

func fqdnAllowWildcard(ctx context.Context, fl validator.FieldLevel) bool {
	val := fl.Field().String()

	if val == "" {
		return false
	}

	return fqdnAllowWildcardRegex.MatchString(val)
}
