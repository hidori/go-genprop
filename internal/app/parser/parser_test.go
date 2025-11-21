package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		fileName      string
		wantError     bool
		wantPackage   string
		errorContains string
	}{
		{
			name:        "valid go file",
			fileName:    "../testdata/valid_syntax_input.go.txt",
			wantError:   false,
			wantPackage: "test",
		},
		{
			name:          "file not found",
			fileName:      "nonexistent.go",
			wantError:     true,
			errorContains: "no such file or directory",
		},
		{
			name:          "success: calls internal parser",
			fileName:      "../testdata/valid_syntax_input.go.txt",
			wantError:     false,
			wantPackage:   "test",
			errorContains: "",
		},
		{
			name:          "invalid syntax test data",
			fileName:      "../testdata/invalid_syntax_input.go.txt",
			wantError:     true,
			errorContains: "",
		},
		{
			name:          "empty filename",
			fileName:      "",
			wantError:     true,
			errorContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := ParseFile(tt.fileName)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, file)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, file)
			assert.Equal(t, tt.wantPackage, file.Name.Name)
			assert.NotEmpty(t, file.Decls)
		})
	}
}
