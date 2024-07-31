package example

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var _validator = validator.New()

func validateFieldValue(v any, tag string) error {
	if err := _validator.Var(v, tag); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type Example2AliasType = int // ignored

type Example2DefinedType int // ignored

type Example2Struct struct {
	value0 int    // ignored
	value1 int    `property:"get"`
	value2 int    `property:"set"`
	value3 int    `property:"get,set"`
	value4 int    `property:"get,set" validate:"required,min=1"`
	id     int    `property:"get"`
	api    string `property:"get"`
	url    string `property:"get"`
	http   string `property:"get"`
}

func NewExample2Struct(v2 int, v3 int, v4 int) (*Example2Struct, error) {
	v := &Example2Struct{}

	v.SetValue2(v3)
	v.SetValue3(v3)

	err := v.SetValue4(v4)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return v, nil
}
