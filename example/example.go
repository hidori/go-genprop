package example

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Struct struct {
	value1 int `property:"get"`
	value2 int `property:"set"`
	value3 int `property:"get,set"`
	value4 int `property:"get,set" validate:"min=1,max=100"`
	value5 int `property:"get,set=private" validate:"min=1,max=100"`
}

var _validator = validator.New()

func validateFieldValue(v any, tag string) error {
	if err := _validator.Var(v, tag); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func NewStruct(v1 int, v2 int, v3 int, v4 int, v5 int) (*Struct, error) {
	v := &Struct{
		value1: v1,
		value2: v2,
		value3: v3,
	}

	err := v.SetValue4(v4)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = v.setValue5(v5)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return v, nil
}
