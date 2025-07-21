package example

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Account represents a user account with validation.
type Account struct {
	id       int    `property:"get"`                                                    // Read-only ID field
	username string `property:"get,set" validate:"required,min=3,max=20"`               // Username with length validation
	email    string `property:"get,set" validate:"required,email"`                      // Email with format validation
	password string `property:"set=private" validate:"required,min=8,containsany=!@#$"` // Password with length and special character validation
}

// validateField validates field values using go-playground/validator.
func validateField(fieldName string, value any, rule string) error {
	if err := validate.Var(value, rule); err != nil {
		return fmt.Errorf("validation failed for field '%s': %s", fieldName, err.Error())
	}

	return nil
}
