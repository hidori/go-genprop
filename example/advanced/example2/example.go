package example2

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type User struct {
	email    string `property:"get,set" validate:"email"`
	password string `property:"set=private" validate:"min=8"`
	score    int    `property:"get,set" validate:"min=0,max=100"`
}

var _validator = validator.New()

func NewUser(email string, password string, score int) (*User, error) {
	user := &User{}

	if err := user.SetEmail(email); err != nil {
		return nil, err
	}

	if err := user.setPassword(password); err != nil {
		return nil, err
	}

	if err := user.SetScore(score); err != nil {
		return nil, err
	}

	return user, nil
}

func validateFieldValue(name string, v any, tag string) error {
	if err := _validator.Var(v, tag); err != nil {
		return errors.Wrapf(err, "validation failed for field '%s'", name)
	}

	return nil
}
