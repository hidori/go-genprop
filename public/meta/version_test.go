package meta

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		pattern      string
		wantNotEmpty bool
	}{
		{
			name:         "success: semantic version format",
			pattern:      `^v\d+\.\d+\.\d+(-[a-zA-Z0-9\-\.]+)?(\+[a-zA-Z0-9\-\.]+)?$`,
			wantNotEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			version := GetVersion()

			if tt.wantNotEmpty {
				require.NotEmpty(t, version)
			}

			if tt.pattern != "" {
				matched, err := regexp.MatchString(tt.pattern, version)
				require.NoError(t, err)
				assert.True(t, matched, "version %q should match pattern %q", version, tt.pattern)
			}
		})
	}
}

func TestGetVersion_Consistency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		callCount int
	}{
		{
			name:      "success: multiple calls return same result",
			callCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			versions := make([]string, tt.callCount)
			for i := 0; i < tt.callCount; i++ {
				versions[i] = GetVersion()
			}

			for i := 1; i < len(versions); i++ {
				assert.Equal(t, versions[0], versions[i])
			}
		})
	}
}
