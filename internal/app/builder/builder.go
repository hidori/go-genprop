package builder

import (
	"flag"
	"fmt"
	"os"
)

// Flags contains all flag pointers.
type Flags struct {
	VersionFlag        *bool
	InitialismFlag     *string
	ValidationFuncFlag *string
	ValidationTagFlag  *string
	NewFunc            *bool
}

// BuildFlagSet creates and configures a flag set.
func BuildFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("genprop", flag.ExitOnError)

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: genprop [flags] <FILE>\n")
		fmt.Fprintf(os.Stderr, "\nA Go code generator that automatically creates getter and setter methods "+
			"for private struct fields based on struct tags.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flagSet.PrintDefaults()
	}

	return flagSet
}

// BuildFlags creates and returns flag pointers from the given flag set.
func BuildFlags(flagSet *flag.FlagSet) *Flags {
	versionFlag := flagSet.Bool("version", false, "show version information")
	initialismFlag := flagSet.String("initialism", "id,url,api", "specify names to which initialism should be applied")
	validationFuncFlag := flagSet.String("validation-func", "validateFieldValue", "specify validation func name")
	validationTagFlag := flagSet.String("validation-tag", "validate", "specify validation tag name")
	newFuncFlag := flagSet.Bool("new-func", false, "generate New function for target struct")

	return &Flags{
		VersionFlag:        versionFlag,
		InitialismFlag:     initialismFlag,
		ValidationFuncFlag: validationFuncFlag,
		ValidationTagFlag:  validationTagFlag,
		NewFunc:            newFuncFlag,
	}
}
