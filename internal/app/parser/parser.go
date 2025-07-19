package parser

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/pkg/errors"
)

// ParseFile parses a Go source file and returns the AST.
func ParseFile(fileName string) (*ast.File, error) {
	file, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}
