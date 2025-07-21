package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		fileName    string
		wantErr     bool
		wantPackage string
	}{
		{
			name:        "valid go file",
			fileName:    "../../../testdata/internal/app/valid_syntax_input.go.txt",
			wantErr:     false,
			wantPackage: "test",
		},
		{
			name:     "file not found",
			fileName: "nonexistent.go",
			wantErr:  true,
		},
		{
			name:        "success: calls internal parser",
			fileName:    "../../../testdata/internal/app/valid_syntax_input.go.txt",
			wantErr:     false,
			wantPackage: "test",
		},
		{
			name:     "invalid syntax test data",
			fileName: "../../../testdata/internal/app/invalid_syntax_input.go.txt",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			fileName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := ParseFile(tt.fileName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, file)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, file)
			assert.Equal(t, tt.wantPackage, file.Name.Name)
			assert.NotEmpty(t, file.Decls)
		})
	}
}
