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
		name            string
		inputFileName   string
		outputFileName  string
		fields          fields
		wantErr         bool
		wantErrContains string
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
			wantErr:         false,
			wantErrContains: "",
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
			wantErr:         true,
			wantErrContains: "invalid tag value",
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
				assert.Contains(t, err.Error(), tt.wantErrContains)

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

// TestFromGenDecl tests edge cases for fromGenDecl method
func TestFromGenDecl(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name     string
		genDecl  *ast.GenDecl
		wantDecl bool
		wantErr  bool
	}{
		{
			name: "success: CONST token returns empty",
			genDecl: &ast.GenDecl{
				Tok: token.CONST,
			},
			wantDecl: false,
			wantErr:  false,
		},
		{
			name: "success: VAR token returns empty",
			genDecl: &ast.GenDecl{
				Tok: token.VAR,
			},
			wantDecl: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decls, err := generator.fromGenDecl(tt.genDecl)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.wantDecl {
				assert.NotEmpty(t, decls)
			} else {
				assert.Empty(t, decls)
			}
		})
	}
}

// TestFromTypeGenDecl tests edge cases for fromTypeGenDecl method
func TestFromTypeGenDecl(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name     string
		genDecl  *ast.GenDecl
		wantDecl bool
		wantErr  bool
	}{
		{
			name: "success: non-TypeSpec returns empty",
			genDecl: &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.ValueSpec{}, // Not a TypeSpec
				},
			},
			wantDecl: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decls, err := generator.fromTypeGenDecl(tt.genDecl)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.wantDecl {
				assert.NotEmpty(t, decls)
			} else {
				assert.Empty(t, decls)
			}
		})
	}
}

// TestFromTypeSpec tests edge cases for fromTypeSpec method
func TestFromTypeSpec(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name     string
		typeSpec *ast.TypeSpec
		wantDecl bool
		wantErr  bool
	}{
		{
			name: "success: non-struct type returns empty",
			typeSpec: &ast.TypeSpec{
				Name: &ast.Ident{Name: "TestType"},
				Type: &ast.InterfaceType{}, // Not a struct
			},
			wantDecl: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decls, err := generator.fromTypeSpec(tt.typeSpec)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.wantDecl {
				assert.NotEmpty(t, decls)
			} else {
				assert.Empty(t, decls)
			}
		})
	}
}

// TestFromField tests edge cases for fromField method
func TestFromField(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name       string
		structName string
		field      *ast.Field
		wantNil    bool
		wantEmpty  bool
		wantErr    bool
	}{
		{
			name:       "success: no tag returns nil",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				// No Tag
			},
			wantNil: true,
		},
		{
			name:       "success: dash tag returns empty",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				Tag:   &ast.BasicLit{Value: "`property:\"-\"`"},
			},
			wantEmpty: true,
		},
		{
			name:       "success: empty tag returns empty",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				Tag:   &ast.BasicLit{Value: "`json:\"test\"`"}, // No property tag
			},
			wantEmpty: true,
		},
		{
			name:       "failure: invalid tag syntax",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				Tag:   &ast.BasicLit{Value: "invalid-quote"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decls, err := generator.fromField(tt.structName, tt.field)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, decls)
			} else if tt.wantEmpty {
				assert.Empty(t, decls)
			} else {
				assert.NotEmpty(t, decls)
			}
		})
	}
}

// TestGetterFuncDecl tests edge cases for getterFuncDecl method
func TestGetterFuncDecl(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name       string
		structName string
		field      *ast.Field
		wantNil    bool
	}{
		{
			name:       "success: anonymous field returns nil",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{}, // No names (anonymous field)
				Type:  &ast.Ident{Name: "string"},
			},
			wantNil: true,
		},
		{
			name:       "success: pointer type returns decl",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "value"}},
				Type:  &ast.StarExpr{X: &ast.Ident{Name: "string"}},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decl := generator.getterFuncDecl(tt.structName, tt.field)

			if tt.wantNil {
				assert.Nil(t, decl)
			} else {
				assert.NotNil(t, decl)
			}
		})
	}
}

// TestSetterFuncDecl tests edge cases for setterFuncDecl method
func TestSetterFuncDecl(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name       string
		prefix     string
		structName string
		field      *ast.Field
		wantNil    bool
	}{
		{
			name:       "success: no tag returns nil",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				// No Tag
			},
			wantNil: true,
		},
		{
			name:       "success: anonymous field returns nil",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{}, // No names
				Type:  &ast.Ident{Name: "string"},
				Tag:   &ast.BasicLit{Value: "`property:\"set\"`"},
			},
			wantNil: true,
		},
		{
			name:       "failure: invalid tag quote returns nil",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				Tag:   &ast.BasicLit{Value: "invalid-quote"}, // Invalid quote syntax
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decl := generator.setterFuncDecl(tt.prefix, tt.structName, tt.field)

			if tt.wantNil {
				assert.Nil(t, decl)
			} else {
				assert.NotNil(t, decl)
			}
		})
	}
}

// TestSetterFuncNoValidationDecl tests edge cases for setterFuncNoValidationDecl method
func TestSetterFuncNoValidationDecl(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name       string
		prefix     string
		structName string
		field      *ast.Field
		wantNil    bool
	}{
		{
			name:       "success: anonymous field returns nil",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{}, // No names
				Type:  &ast.Ident{Name: "string"},
			},
			wantNil: true,
		},
		{
			name:       "success: slice type returns decl",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "items"}},
				Type:  &ast.ArrayType{Elt: &ast.Ident{Name: "string"}},
				Tag:   &ast.BasicLit{Value: "`property:\"set\"`"},
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decl := generator.setterFuncNoValidationDecl(tt.prefix, tt.structName, tt.field)

			if tt.wantNil {
				assert.Nil(t, decl)
			} else {
				assert.NotNil(t, decl)
			}
		})
	}
}

// TestSetterFuncWithValidationDecl tests edge cases for setterFuncWithValidationDecl method
func TestSetterFuncWithValidationDecl(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name          string
		prefix        string
		structName    string
		field         *ast.Field
		validationTag string
		wantNil       bool
	}{
		{
			name:       "success: no tag returns nil",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				// No Tag
			},
			validationTag: "required",
			wantNil:       true,
		},
		{
			name:       "success: anonymous field returns nil",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{}, // No names
				Type:  &ast.Ident{Name: "string"},
				Tag:   &ast.BasicLit{Value: "`property:\"set\" validate:\"required\"`"},
			},
			validationTag: "required",
			wantNil:       true,
		},
		{
			name:       "success: no validation tag returns nil",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				Tag:   &ast.BasicLit{Value: "`property:\"set\"`"}, // No validate tag
			},
			validationTag: "",
			wantNil:       true,
		},
		{
			name:       "failure: invalid tag quote returns nil",
			prefix:     "Set",
			structName: "TestStruct",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
				Tag:   &ast.BasicLit{Value: "invalid-quote"}, // Invalid quote syntax
			},
			validationTag: "required",
			wantNil:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decl := generator.setterFuncWithValidationDecl(tt.prefix, tt.structName, tt.field, tt.validationTag)

			if tt.wantNil {
				assert.Nil(t, decl)
			} else {
				assert.NotNil(t, decl)
			}
		})
	}
}

// TestPrepareFieldName tests edge cases for prepareFieldName method
func TestPrepareFieldName(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name      string
		fieldName string
		want      string
	}{
		{
			name:      "success: api initialism",
			fieldName: "apiKey",
			want:      "APIKey",
		},
		{
			name:      "success: id initialism",
			fieldName: "idValue",
			want:      "IDValue",
		},
		{
			name:      "success: no lowercase prefix",
			fieldName: "TestField",
			want:      "TestField",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := generator.prepareFieldName(tt.fieldName)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestBuildSetterFuncType tests edge cases for buildSetterFuncType method
func TestBuildSetterFuncType(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name      string
		field     *ast.Field
		withError bool
		wantNil   bool
	}{
		{
			name: "success: with error",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "value"}},
				Type:  &ast.Ident{Name: "string"},
			},
			withError: true,
			wantNil:   false,
		},
		{
			name: "success: without error",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "value"}},
				Type:  &ast.Ident{Name: "string"},
			},
			withError: false,
			wantNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			funcType := generator.buildSetterFuncType(tt.field, tt.withError)

			if tt.wantNil {
				assert.Nil(t, funcType)
			} else {
				require.NotNil(t, funcType)
				if tt.withError {
					assert.NotNil(t, funcType.Results)
				} else {
					assert.Nil(t, funcType.Results)
				}
			}
		})
	}
}

func TestGenerator_processDirective(t *testing.T) {
	t.Parallel()

	generator := NewGenerator(&GeneratorConfig{
		TagName: tagName,
	})

	field := &ast.Field{
		Names: []*ast.Ident{{Name: "testField"}},
		Type:  &ast.Ident{Name: "string"},
		Tag:   &ast.BasicLit{Kind: token.STRING, Value: "`property:\"get\"`"},
	}

	tests := []struct {
		name      string
		directive string
		wantErr   bool
	}{
		{
			name:      "success: get directive",
			directive: "get",
			wantErr:   false,
		},
		{
			name:      "success: set directive",
			directive: "set",
			wantErr:   false,
		},
		{
			name:      "success: set=private directive",
			directive: "set=private",
			wantErr:   false,
		},
		{
			name:      "failure: invalid directive",
			directive: "invalid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decl, err := generator.processDirective(tt.directive, "TestStruct", field)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, decl)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, decl)
			}
		})
	}
}

func TestFromFieldList(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name       string
		structName string
		fieldList  *ast.FieldList
		wantEmpty  bool
		wantErr    bool
	}{
		{
			name:       "success: empty field list returns empty",
			structName: "TestStruct",
			fieldList:  &ast.FieldList{List: []*ast.Field{}},
			wantEmpty:  true,
		},
		{
			name:       "success: field with property tag returns decl",
			structName: "TestStruct",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "TestField"}},
						Type:  &ast.Ident{Name: "string"},
						Tag:   &ast.BasicLit{Value: "`property:\"get\"`"},
					},
				},
			},
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decls, err := generator.fromFieldList(tt.structName, tt.fieldList)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.wantEmpty {
				assert.Empty(t, decls)
			} else {
				assert.NotEmpty(t, decls)
			}
		})
	}
}

func TestBuildRecvFieldList(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name       string
		structName string
		wantNil    bool
	}{
		{
			name:       "success: returns receiver field list",
			structName: "TestStruct",
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fieldList := generator.buildRecvFieldList(tt.structName)

			if tt.wantNil {
				assert.Nil(t, fieldList)
			} else {
				assert.NotNil(t, fieldList)
				assert.NotEmpty(t, fieldList.List)
			}
		})
	}
}

func TestBuildValidationBody(t *testing.T) {
	t.Parallel()

	config := &GeneratorConfig{
		TagName:        "property",
		Initialism:     []string{"api", "id"},
		ValidationFunc: "validate",
		ValidationTag:  "validate",
	}
	generator := NewGenerator(config)

	tests := []struct {
		name    string
		field   *ast.Field
		tag     string
		wantNil bool
	}{
		{
			name: "success: returns validation body",
			field: &ast.Field{
				Names: []*ast.Ident{{Name: "TestField"}},
				Type:  &ast.Ident{Name: "string"},
			},
			tag:     "required",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body := generator.buildValidationBody(tt.field, tt.tag)

			if tt.wantNil {
				assert.Nil(t, body)
			} else {
				assert.NotNil(t, body)
				assert.NotEmpty(t, body.List)
			}
		})
	}
}
