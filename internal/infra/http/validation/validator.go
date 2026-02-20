package validation

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var validatorOnce sync.Once
var validate *validator.Validate

func GetValidator() *validator.Validate {
	validatorOnce.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
	})
	return validate
}
