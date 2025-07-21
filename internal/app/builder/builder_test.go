package builder

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFlagSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		assert func(t *testing.T, flagSet *flag.FlagSet)
	}{
		{
			name: "success: creates flag set with correct name",
			assert: func(t *testing.T, flagSet *flag.FlagSet) {
				assert.Equal(t, "genprop", flagSet.Name())
			},
		},
		{
			name: "success: configures usage function",
			assert: func(t *testing.T, flagSet *flag.FlagSet) {
				assert.NotNil(t, flagSet.Usage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			flagSet := BuildFlagSet()

			// Common assertions
			require.NotNil(t, flagSet)

			// Test case specific assertions
			if tt.assert != nil {
				tt.assert(t, flagSet)
			}
		})
	}
}

func TestBuildFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		args               []string
		wantInitialism     string
		wantValidationFunc string
		wantValidationTag  string
		wantVersion        bool
	}{
		{
			name:               "success: creates flags with default values",
			args:               []string{},
			wantInitialism:     "id,url,api",
			wantValidationFunc: "validateFieldValue",
			wantValidationTag:  "validate",
			wantVersion:        false,
		},
		{
			name: "success: respects parsed flag values",
			args: []string{
				"-initialism", "custom,id",
				"-validation-func", "customValidate",
				"-validation-tag", "customTag",
				"-version",
			},
			wantInitialism:     "custom,id",
			wantValidationFunc: "customValidate",
			wantValidationTag:  "customTag",
			wantVersion:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			flags := BuildFlags(flagSet)

			require.NotNil(t, flags)
			require.NotNil(t, flags.InitialismFlag)
			require.NotNil(t, flags.ValidationFuncFlag)
			require.NotNil(t, flags.ValidationTagFlag)
			require.NotNil(t, flags.VersionFlag)

			if len(tt.args) > 0 {
				err := flagSet.Parse(tt.args)
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantInitialism, *flags.InitialismFlag)
			assert.Equal(t, tt.wantValidationFunc, *flags.ValidationFuncFlag)
			assert.Equal(t, tt.wantValidationTag, *flags.ValidationTagFlag)
			assert.Equal(t, tt.wantVersion, *flags.VersionFlag)
		})
	}
}

func TestFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "success: flags struct contains all required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			flags := BuildFlags(flagSet)

			// Test that all fields are properly initialized
			assert.IsType(t, &Flags{}, flags)
			assert.IsType(t, (*string)(nil), flags.InitialismFlag)
			assert.IsType(t, (*string)(nil), flags.ValidationFuncFlag)
			assert.IsType(t, (*string)(nil), flags.ValidationTagFlag)
			assert.IsType(t, (*bool)(nil), flags.VersionFlag)
		})
	}
}
