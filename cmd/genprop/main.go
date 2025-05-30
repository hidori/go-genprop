package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/hidori/go-genprop/generator"
	"github.com/hidori/go-genprop/meta"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

const (
	doNotEdit = "// Code generated by github.com/hidori/go-genprop/cmd/genprop DO NOT EDIT."
	tagName   = "property"
)

var (
	initialismFlag    = flag.String("initialism", "id,url,api", "specify names to which initialism should be applied")
	validatioFuncFlag = flag.String("validation-func", "validateFieldValue", "specify validation func name")
	validatioTagFlag  = flag.String("validation-tag", "validate", "specify validation tag name")
)

func main() {
	name := path.Base(os.Args[0])

	if slices.Contains(os.Args, "-version") {
		fmt.Printf("%s %s\n", name, meta.GetVersion())

		return
	}

	flag.Usage = func() {
		fmt.Printf("Usage:\n  %s [OPTION]... <FILE>\n\nOption(s):\n", name)
		fmt.Println("  -version\n        show version information")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()

		return
	}

	err := generate(os.Stdout, args[0])
	if err != nil {
		log.Fatal(err)
	}
}

func generate(writer io.Writer, fileName string) error {
	file, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
	if err != nil {
		return errors.Wrap(err, "fail to parser.ParseFile()")
	}

	generator := generator.NewGenerator(&generator.GeneratorConfig{
		TagName:        tagName,
		Initialism:     strings.Split(*initialismFlag, ","),
		ValidationFunc: *validatioFuncFlag,
		ValidationTag:  *validatioTagFlag,
	})

	decls, err := generator.Generate(token.NewFileSet(), file)
	if err != nil {
		return errors.Wrap(err, "fail to generator.Generate()")
	}

	buffer := bytes.NewBuffer([]byte{})

	err = format.Node(buffer, token.NewFileSet(), &ast.File{
		Name:  ast.NewIdent(file.Name.Name),
		Decls: decls,
	})
	if err != nil {
		return errors.Wrap(err, "fail to format.Node()")
	}

	cooked, err := imports.Process("", buffer.Bytes(), &imports.Options{FormatOnly: false})
	if err != nil {
		return errors.Wrap(err, "fail to imports.Process()")
	}

	_, _ = fmt.Fprintln(writer, doNotEdit)
	_, _ = writer.Write(cooked)

	return nil
}
