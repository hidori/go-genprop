package genprop

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tagName = "property"

func TestGenerator_Generate(t *testing.T) {
	type fields struct {
		config *GeneratorConfig
	}
	tests := []struct {
		name           string
		input          string
		output         string
		fields         fields
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:   "success: returns []ast.Decl",
			input:  "./test/data/success_input.go",
			output: "./test/data/success_output.txt",
			fields: fields{
				config: &GeneratorConfig{
					TagName:    tagName,
					Initialism: []string{"api"},
				},
			},
		},
		{
			name:  "fail: returns []ast.Decl",
			input: "./test/data/fail_input.go",
			fields: fields{
				config: &GeneratorConfig{
					TagName:    tagName,
					Initialism: []string{"api"},
				},
			},
			wantErr:        true,
			wantErrMessage: "invalid tag value 'undefined'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()

			f, err := parser.ParseFile(token.NewFileSet(), tt.input, nil, parser.AllErrors)
			if err != nil {
				t.Errorf("fail to parser.ParseFile() tt.input=%v", tt.input)
				return
			}

			got, err := NewGenerator(tt.fields.config).Generate(fset, f)
			if err != nil && tt.wantErr {
				assert.Contains(t, err.Error(), tt.wantErrMessage)
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			_want := bytes.NewBuffer([]byte{})
			{
				f, err := parser.ParseFile(token.NewFileSet(), tt.output, nil, parser.AllErrors)
				if err != nil {
					t.Errorf("fail to parser.ParseFile() tt.output=%v", tt.input)
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
