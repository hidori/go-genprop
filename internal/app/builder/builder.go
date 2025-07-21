package builder

import (
	"flag"
	"fmt"
	"os"
)

// Flags contains all flag pointers.
type Flags struct {
	GenerateNewFunc    *bool
	InitialismFlag     *string
	ValidationFuncFlag *string
	ValidationTagFlag  *string
	VersionFlag        *bool
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
	generateNewFuncFlag := flagSet.Bool("generate-new-func", false, "generate New function for target struct")
	initialismFlag := flagSet.String("initialism", "id,url,api", "specify names to which initialism should be applied")
	validationFuncFlag := flagSet.String("validation-func", "validateFieldValue", "specify validation func name")
	validationTagFlag := flagSet.String("validation-tag", "validate", "specify validation tag name")
	versionFlag := flagSet.Bool("version", false, "show version information")

	return &Flags{
		GenerateNewFunc:    generateNewFuncFlag,
		InitialismFlag:     initialismFlag,
		ValidationFuncFlag: validationFuncFlag,
		ValidationTagFlag:  validationTagFlag,
		VersionFlag:        versionFlag,
	}
}
