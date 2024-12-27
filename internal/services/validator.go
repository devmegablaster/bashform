package services

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		v: validator.New(),
	}
}

// Validate a struct using the validator
func (v *Validator) Validate(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		errors := v.parseError(err)
		errStr := ""
		for field, err := range errors {
			errStr += field + ": " + err + ", "
		}

		return fmt.Errorf("Validation errors: %s", errStr)
	}
	return nil
}

// parseError takes a validation error and returns a map of field names to error messages
func (v *Validator) parseError(err error) map[string]string {
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = err.Tag()
	}
	return errors
}
