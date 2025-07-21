package generator

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/hidori/go-genprop/public/generator"
	"github.com/pkg/errors"
)

const (
	tagName = "property"
)

// GenerateCode generates AST declarations for getter and setter methods based on the given file and configuration.
func GenerateCode(
	file *ast.File,
	initialismFlag, validationFuncFlag, validationTagFlag string,
	generateNewFunc bool,
) ([]ast.Decl, error) {
	generator := generator.NewGenerator(&generator.GeneratorConfig{
		TagName:         tagName,
		GenerateNewFunc: generateNewFunc,
		Initialism:      strings.Split(initialismFlag, ","),
		ValidationFunc:  validationFuncFlag,
		ValidationTag:   validationTagFlag,
	})

	decls, err := generator.Generate(token.NewFileSet(), file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return decls, nil
}
