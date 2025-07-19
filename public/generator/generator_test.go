package generator

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tagName = "property"

func TestGenerator_Generate(t *testing.T) {
	t.Parallel()

	type fields struct {
		config *GeneratorConfig
	}

	tests := []struct {
		name           string
		inputFileName  string
		outputFileName string
		fields         fields
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:           "success: returns ast.Decl",
			inputFileName:  "../../testdata/public/generator/basic_getset_input.go.txt",
			outputFileName: "../../testdata/public/generator/basic_getset_output.txt",
			fields: fields{
				config: &GeneratorConfig{
					TagName:    tagName,
					Initialism: []string{"api"},
				},
			},
		},
		{
			name:          "failure: returns error for invalid tag",
			inputFileName: "../../testdata/public/generator/invalid_directive_input.go.txt",
			fields: fields{
				config: &GeneratorConfig{
					TagName:    tagName,
					Initialism: []string{"api"},
				},
			},
			wantErr:        true,
			wantErrMessage: "invalid tag value",
		},
		{
			name:           "success: returns ast.Decl with validation",
			inputFileName:  "../../testdata/public/generator/validation_input.go.txt",
			outputFileName: "../../testdata/public/generator/validation_output.txt",
			fields: fields{
				config: &GeneratorConfig{
					TagName:        tagName,
					Initialism:     []string{"api"},
					ValidationFunc: "validateFieldValue",
					ValidationTag:  "validate",
				},
			},
		},
		{
			name:           "success: returns ast.Decl with private setter",
			inputFileName:  "../../testdata/public/generator/private_setter_input.go.txt",
			outputFileName: "../../testdata/public/generator/private_setter_output.txt",
			fields: fields{
				config: &GeneratorConfig{
					TagName:    tagName,
					Initialism: []string{"api"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fset := token.NewFileSet()

			f, err := parser.ParseFile(token.NewFileSet(), tt.inputFileName, nil, parser.AllErrors)
			if err != nil {
				t.Errorf("fail to parser.ParseFile() tt.inputFileName=%v", tt.inputFileName)

				return
			}

			got, err := NewGenerator(tt.fields.config).Generate(fset, f)
			if err != nil && tt.wantErr {
				assert.Contains(t, err.Error(), tt.wantErrMessage)

				return
			}

			require.NoError(t, err)

			_want := bytes.NewBuffer([]byte{})

			{
				f, err := parser.ParseFile(token.NewFileSet(), tt.outputFileName, nil, parser.AllErrors)
				if err != nil {
					t.Errorf("fail to parser.ParseFile() tt.outputFileName=%v", tt.outputFileName)

					return
				}

				format.Node(_want, fset, f.Decls)
			}

			_got := bytes.NewBuffer([]byte{})

			format.Node(_got, fset, got)

			if !assert.Equal(t, _want.String(), _got.String()) {
				return
			}
		})
	}
}

// TestCoverageEdgeCases tests edge cases to improve coverage
func TestCoverageEdgeCases(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	t.Run("success: fromGenDecl with CONST token", func(t *testing.T) {
		t.Parallel()

		genDecl := &ast.GenDecl{
			Tok: token.CONST,
		}

		decls, err := generator.fromGenDecl(genDecl)
		require.NoError(t, err)
		assert.Empty(t, decls)
	})

	t.Run("success: fromGenDecl with VAR token", func(t *testing.T) {
		t.Parallel()

		genDecl := &ast.GenDecl{
			Tok: token.VAR,
		}

		decls, err := generator.fromGenDecl(genDecl)
		require.NoError(t, err)
		assert.Empty(t, decls)
	})

	t.Run("success: fromTypeGenDecl with non-TypeSpec", func(t *testing.T) {
		t.Parallel()

		genDecl := &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.ValueSpec{}, // Not a TypeSpec
			},
		}

		decls, err := generator.fromTypeGenDecl(genDecl)
		require.NoError(t, err)
		assert.Empty(t, decls)
	})

	t.Run("success: fromTypeSpec with non-struct type", func(t *testing.T) {
		t.Parallel()

		typeSpec := &ast.TypeSpec{
			Name: &ast.Ident{Name: "TestType"},
			Type: &ast.InterfaceType{}, // Not a struct
		}

		decls, err := generator.fromTypeSpec(typeSpec)
		require.NoError(t, err)
		assert.Empty(t, decls)
	})

	t.Run("success: fromField with no tag", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			// No Tag
		}

		decls, err := generator.fromField("TestStruct", field)
		require.NoError(t, err)
		assert.Nil(t, decls)
	})

	t.Run("success: fromField with dash tag", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			Tag:   &ast.BasicLit{Value: "`property:\"-\"`"},
		}

		decls, err := generator.fromField("TestStruct", field)
		require.NoError(t, err)
		assert.Empty(t, decls)
	})

	t.Run("success: fromField with empty tag", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			Tag:   &ast.BasicLit{Value: "`json:\"test\"`"}, // No property tag
		}

		decls, err := generator.fromField("TestStruct", field)
		require.NoError(t, err)
		assert.Empty(t, decls)
	})

	t.Run("success: getterFuncDecl with anonymous field", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{}, // No names (anonymous field)
			Type:  &ast.Ident{Name: "string"},
		}

		decl := generator.getterFuncDecl("TestStruct", field)
		assert.Nil(t, decl)
	})

	t.Run("success: setterFuncDecl with no tag", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			// No Tag
		}

		decl := generator.setterFuncDecl("Set", "TestStruct", field)
		assert.Nil(t, decl)
	})

	t.Run("success: setterFuncDecl with anonymous field", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{}, // No names
			Type:  &ast.Ident{Name: "string"},
			Tag:   &ast.BasicLit{Value: "`property:\"set\"`"},
		}

		decl := generator.setterFuncDecl("Set", "TestStruct", field)
		assert.Nil(t, decl)
	})

	t.Run("success: setterFuncNoValidationDecl with anonymous field", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{}, // No names
			Type:  &ast.Ident{Name: "string"},
		}

		decl := generator.setterFuncNoValidationDecl("Set", "TestStruct", field)
		assert.Nil(t, decl)
	})

	t.Run("success: setterFuncWithValidationDecl with no tag", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			// No Tag
		}

		decl := generator.setterFuncWithValidationDecl("Set", "TestStruct", field, "required")
		assert.Nil(t, decl)
	})

	t.Run("success: setterFuncWithValidationDecl with anonymous field", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{}, // No names
			Type:  &ast.Ident{Name: "string"},
			Tag:   &ast.BasicLit{Value: "`property:\"set\" validate:\"required\"`"},
		}

		decl := generator.setterFuncWithValidationDecl("Set", "TestStruct", field, "required")
		assert.Nil(t, decl)
	})

	t.Run("success: setterFuncWithValidationDecl with no validation tag", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			Tag:   &ast.BasicLit{Value: "`property:\"set\"`"}, // No validate tag
		}

		decl := generator.setterFuncWithValidationDecl("Set", "TestStruct", field, "")
		assert.Nil(t, decl)
	})

	t.Run("failure: fromField with invalid tag syntax", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			Tag:   &ast.BasicLit{Value: "invalid-tag-syntax"}, // Invalid syntax
		}

		_, err := generator.fromField("TestStruct", field)
		assert.Error(t, err)
	})

	t.Run("success: prepareFieldName with initialism", func(t *testing.T) {
		t.Parallel()

		result := generator.prepareFieldName("apiKey")
		assert.Equal(t, "APIKey", result)
	})

	t.Run("success: prepareFieldName with id initialism", func(t *testing.T) {
		t.Parallel()

		result := generator.prepareFieldName("idValue")
		assert.Equal(t, "IDValue", result)
	})

	t.Run("success: prepareFieldName with no lowercase prefix", func(t *testing.T) {
		t.Parallel()

		result := generator.prepareFieldName("TestField")
		assert.Equal(t, "TestField", result)
	})
}

// TestCoverageGeneratorFailures tests error conditions
func TestCoverageGeneratorFailures(t *testing.T) {
	t.Parallel()

	t.Run("failure: setterFuncDecl with invalid tag quote", func(t *testing.T) {
		t.Parallel()

		config := &GeneratorConfig{
			TagName:        "property",
			ValidationFunc: "validate",
			ValidationTag:  "validate",
		}
		generator := NewGenerator(config)

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			Tag:   &ast.BasicLit{Value: "invalid-quote"}, // Invalid quote syntax
		}

		decl := generator.setterFuncDecl("Set", "TestStruct", field)
		assert.Nil(t, decl)
	})

	t.Run("failure: setterFuncWithValidationDecl with invalid tag quote", func(t *testing.T) {
		t.Parallel()

		config := &GeneratorConfig{
			TagName:        "property",
			ValidationFunc: "validate",
			ValidationTag:  "validate",
		}
		generator := NewGenerator(config)

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "TestField"}},
			Type:  &ast.Ident{Name: "string"},
			Tag:   &ast.BasicLit{Value: "invalid-quote"}, // Invalid quote syntax
		}

		decl := generator.setterFuncWithValidationDecl("Set", "TestStruct", field, "required")
		assert.Nil(t, decl)
	})
}

// TestCoverageComplexFieldTypes tests with complex field types
func TestCoverageComplexFieldTypes(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	t.Run("success: getter with pointer type", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "value"}},
			Type:  &ast.StarExpr{X: &ast.Ident{Name: "string"}},
		}

		decl := generator.getterFuncDecl("TestStruct", field)
		assert.NotNil(t, decl)
	})

	t.Run("success: setter with slice type", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "items"}},
			Type:  &ast.ArrayType{Elt: &ast.Ident{Name: "string"}},
			Tag:   &ast.BasicLit{Value: "`property:\"set\"`"},
		}

		decl := generator.setterFuncNoValidationDecl("Set", "TestStruct", field)
		assert.NotNil(t, decl)
	})

	t.Run("success: buildSetterFuncType with error", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "value"}},
			Type:  &ast.Ident{Name: "string"},
		}

		funcType := generator.buildSetterFuncType(field, true)
		assert.NotNil(t, funcType)
		assert.NotNil(t, funcType.Results)
	})

	t.Run("success: buildSetterFuncType without error", func(t *testing.T) {
		t.Parallel()

		field := &ast.Field{
			Names: []*ast.Ident{{Name: "value"}},
			Type:  &ast.Ident{Name: "string"},
		}

		funcType := generator.buildSetterFuncType(field, false)
		assert.NotNil(t, funcType)
		assert.Nil(t, funcType.Results)
	})
}
