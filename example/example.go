package example

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Struct struct {
	value1 int `property:"get"`
	value2 int `property:"set"`
	value3 int `property:"get,set"`
	value4 int `property:"set=private"`
	value5 int `property:"get,set" validate:"min=1,max=100"`
	value6 int `property:"get,set=private" validate:"min=1,max=100"`
}

var _validator = validator.New()

func validateFieldValue(name string, v any, tag string) error {
	if err := _validator.Var(v, tag); err != nil {
		return errors.Wrapf(errors.WithStack(err), "fail to validator.Var() name='%s'", name)
	}

	return nil
}

func NewStruct(v1 int, v2 int, v3 int, v4 int, v5 int, v6 int) (*Struct, error) {
	v := &Struct{
		value1: v1, // has no setter
	}

	v.SetValue2(v2)
	v.SetValue3(v3)
	v.setValue4(v4)

	err := v.SetValue5(v5)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = v.setValue6(v6)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return v, nil
}
