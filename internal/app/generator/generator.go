package generator

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/hidori/go-genprop/public/generator"
	"github.com/pkg/errors"
)

const (
	propertyTag = "property"
)

// GenerateCode generates AST declarations for getter and setter methods based on the given file and configuration.
func GenerateCode(
	file *ast.File,
	initialism, validationFunc, validationTag string,
	newFunc bool,
) ([]ast.Decl, error) {
	generator := generator.NewGenerator(&generator.GeneratorConfig{
		PropertyTag:    propertyTag,
		Initialism:     strings.Split(initialism, ","),
		ValidationFunc: validationFunc,
		ValidationTag:  validationTag,
		NewFunc:        newFunc,
	})

	decls, err := generator.Generate(token.NewFileSet(), file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return decls, nil
}
