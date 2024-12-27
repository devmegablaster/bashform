package services

import (
	"github.com/devmegablaster/bashform/internal/types"
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
func (v *Validator) Validate(i interface{}) types.ServiceErrors {
	if err := v.v.Struct(i); err != nil {
		return v.parseError(err)
	}
	return nil
}

// parseError takes a validation error and returns a map of field names to error messages
func (v *Validator) parseError(err error) types.ServiceErrors {
	errors := types.ServiceErrors{}
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = err.Tag()
	}
	return errors
}
