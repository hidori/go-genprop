package advanced

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// User represents a user with validation-enabled fields.
type User struct {
	id    int    `property:"get"`                                       // Read-only ID field
	name  string `property:"get,set"`                                   // Name with both getter and setter
	email string `property:"get,set=private" validate:"required,email"` // Email with private setter and validation
}

var _validator = validator.New()

// NewUser creates a new User instance with validation.
func NewUser(id int, name, email string) (*User, error) {
	user := &User{id: id}
	user.SetName(name)

	if err := user.setEmail(email); err != nil {
		return nil, err
	}

	return user, nil
}

// validateFieldValue validates a field value using the specified validation tag.
func validateFieldValue(name string, v any, tag string) error {

	if err := _validator.Var(v, tag); err != nil {

		return errors.Wrapf(err, "validation failed for field '%s'", name)
	}

	return nil
}
