package generator

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		file    *ast.File
		wantErr bool
	}{
		{
			name: "success: calls internal generator",
			file: &ast.File{
				Name:  ast.NewIdent("test"),
				Decls: []ast.Decl{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decls, err := GenerateCode(tt.file, "id,url,api", "validateFieldValue", "validate", false)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, decls)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
