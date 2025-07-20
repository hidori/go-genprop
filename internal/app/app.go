package app

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hidori/go-genprop/internal/app/formatter"
	"github.com/hidori/go-genprop/internal/app/generator"
	"github.com/hidori/go-genprop/internal/app/parser"
	"github.com/hidori/go-genprop/public/meta"
	"github.com/pkg/errors"
)

// Run executes the CLI application with command line arguments.
func Run(args []string) error {
	flagSet := flag.NewFlagSet("genprop", flag.ExitOnError)
	initialismFlagFS := flagSet.String("initialism", "id,url,api", "specify names to which initialism should be applied")
	validationFuncFlagFS := flagSet.String("validation-func", "validateFieldValue", "specify validation func name")
	validationTagFlagFS := flagSet.String("validation-tag", "validate", "specify validation tag name")
	versionFlagFS := flagSet.Bool("version", false, "show version information")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: genprop [flags] <FILE>\n")
		fmt.Fprintf(os.Stderr, "\nA Go code generator that automatically creates getter and setter methods "+
			"for private struct fields based on struct tags.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flagSet.PrintDefaults()
	}

	err := flagSet.Parse(args)
	if err != nil {
		return errors.WithStack(err)
	}

	if *versionFlagFS {
		fmt.Println(meta.GetVersion())
		return nil
	}

	parsedArgs := flagSet.Args()

	if len(parsedArgs) == 0 {
		flagSet.Usage()

		return errors.New("file argument is required")
	}

	if len(parsedArgs) != 1 {
		flagSet.Usage()

		return errors.New("exactly one file argument is required")
	}

	return generate(os.Stdout, parsedArgs[0], *initialismFlagFS, *validationFuncFlagFS, *validationTagFlagFS)
}

func generate(writer io.Writer, fileName string, initialismFlag, validationFuncFlag, validationTagFlag string) error {
	file, err := parser.ParseFile(fileName)
	if err != nil {
		return errors.Wrap(err, "failed to parse file")
	}

	decls, err := generator.GenerateCode(file, initialismFlag, validationFuncFlag, validationTagFlag)
	if err != nil {
		return errors.Wrap(err, "failed to generate code")
	}

	err = formatter.WriteOutput(writer, file.Name.Name, decls)
	if err != nil {
		return errors.Wrap(err, "failed to write output")
	}

	return nil
}
