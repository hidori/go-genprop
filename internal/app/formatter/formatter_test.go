package formatter

import (
	"bytes"
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"github.com/hidori/go-astutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		packageName  string
		decls        []ast.Decl
		wantContains []string
		wantError    bool
		assert       func(t *testing.T, output string)
	}{
		{
			name:        "simple function",
			packageName: "test",
			decls: []ast.Decl{
				&ast.FuncDecl{
					Name: astutil.NewIdent("TestFunc"),
					Type: astutil.NewFuncType(nil, nil, nil),
					Body: astutil.NewBlockStmt([]ast.Stmt{
						astutil.NewReturnStmt(nil),
					}),
				},
			},
			wantContains: []string{
				"// Code generated by github.com/hidori/go-genprop/cmd/genprop DO NOT EDIT.",
				"package test",
				"func TestFunc()",
			},
			assert: func(t *testing.T, output string) {
				// Basic assertions only
			},
		},
		{
			name:        "empty declarations",
			packageName: "empty",
			decls:       []ast.Decl{},
			wantContains: []string{
				"// Code generated by github.com/hidori/go-genprop/cmd/genprop DO NOT EDIT.",
				"package empty",
			},
			assert: func(t *testing.T, output string) {
				// Basic assertions only
			},
		},
		{
			name:        "function with fmt import",
			packageName: "withimports",
			decls: []ast.Decl{
				&ast.FuncDecl{
					Name: astutil.NewIdent("PrintHello"),
					Type: astutil.NewFuncType(nil, nil, nil),
					Body: astutil.NewBlockStmt([]ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   astutil.NewIdent("fmt"),
									Sel: astutil.NewIdent("Println"),
								},
								Args: []ast.Expr{
									astutil.NewBasicLit(token.STRING, `"hello"`),
								},
							},
						},
					}),
				},
			},
			wantContains: []string{
				"package withimports",
				"fmt.Println",
				"import",
			},
			assert: func(t *testing.T, output string) {
				// Check if fmt import is actually added
				lines := strings.Split(output, "\n")
				var hasImport bool
				for _, line := range lines {
					if strings.Contains(line, `"fmt"`) {
						hasImport = true
						break
					}
				}
				assert.True(t, hasImport, "Expected fmt import to be added automatically")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := WriteOutput(&buf, tt.packageName, tt.decls)

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			output := buf.String()

			for _, want := range tt.wantContains {
				assert.Contains(t, output, want)
			}

			// Execute test case specific assertions
			if tt.assert != nil {
				tt.assert(t, output)
			}
		})
	}
}
