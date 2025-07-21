package app

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hidori/go-genprop/internal/app/builder"
	"github.com/hidori/go-genprop/internal/app/formatter"
	"github.com/hidori/go-genprop/internal/app/generator"
	"github.com/hidori/go-genprop/internal/app/parser"
	"github.com/hidori/go-genprop/public/meta"
	"github.com/pkg/errors"
)

// Run executes the CLI application with command line arguments.
func Run(args []string) error {
	flagSet := builder.BuildFlagSet()
	flags := builder.BuildFlags(flagSet)

	err := flagSet.Parse(args)
	if err != nil {
		return errors.WithStack(err)
	}

	if *flags.VersionFlag {
		fmt.Println(meta.GetVersion())
		return nil
	}

	fileName, err := getFileName(flagSet)
	if err != nil {
		return errors.WithStack(err)
	}

	return generate(
		os.Stdout,
		fileName,
		*flags.NewFunc,
		*flags.InitialismFlag,
		*flags.ValidationFuncFlag,
		*flags.ValidationTagFlag,
	)
}

func getFileName(flagSet *flag.FlagSet) (string, error) {
	parsedArgs := flagSet.Args()

	if len(parsedArgs) == 0 {
		flagSet.Usage()

		return "", errors.New("file argument is required")
	}

	if len(parsedArgs) != 1 {
		flagSet.Usage()

		return "", errors.New("exactly one file argument is required")
	}

	return parsedArgs[0], nil
}

func generate(
	writer io.Writer,
	fileName string,
	newFunc bool,
	initialism, validationFunc, validationTag string,
) error {
	file, err := parser.ParseFile(fileName)
	if err != nil {
		return errors.Wrap(err, "failed to parse file")
	}

	decls, err := generator.GenerateCode(file, initialism, validationFunc, validationTag, newFunc)
	if err != nil {
		return errors.Wrap(err, "failed to generate code")
	}

	err = formatter.WriteOutput(writer, file.Name.Name, decls)
	if err != nil {
		return errors.Wrap(err, "failed to write output")
	}

	return nil
}
