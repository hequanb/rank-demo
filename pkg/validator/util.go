package validator

import "github.com/go-playground/validator/v10"

func IsValidationErrors(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(validator.ValidationErrors)
	return ok
}
