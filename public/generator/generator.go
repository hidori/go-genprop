package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hidori/go-astutil"
	"github.com/hidori/go-typeutil"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GeneratorConfig holds configuration for the code generator.
type GeneratorConfig struct {
	TagName         string
	GenerateNewFunc bool
	Initialism      []string
	ValidationFunc  string
	ValidationTag   string
}

// Generator generates getter and setter methods for struct fields.
type Generator struct {
	config *GeneratorConfig
}

// NewGenerator creates a new Generator with the given configuration.
func NewGenerator(config *GeneratorConfig) *Generator {
	return &Generator{
		config: config,
	}
}

// Generate generates getter and setter methods for struct fields based on tags.
func (g *Generator) Generate(fileSet *token.FileSet, file *ast.File) ([]ast.Decl, error) {
	var decls []ast.Decl

	for _, d := range file.Decls {
		genDecl := typeutil.AsOrEmpty[*ast.GenDecl](d)

		if genDecl == nil {
			continue
		}

		_decls, err := g.fromGenDecl(genDecl)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		decls = append(decls, _decls...)
	}

	return decls, nil
}

//nolint:exhaustive // Only IMPORT/TYPE tokens are processed, others handled by default case
func (g *Generator) fromGenDecl(genDecl *ast.GenDecl) ([]ast.Decl, error) {
	switch genDecl.Tok {
	case token.IMPORT:
		return []ast.Decl{genDecl}, nil

	case token.TYPE:
		return g.fromTypeGenDecl(genDecl)

	default:
		return []ast.Decl{}, nil
	}
}

func (g *Generator) fromTypeGenDecl(genDecl *ast.GenDecl) ([]ast.Decl, error) {
	var decls []ast.Decl

	for _, s := range genDecl.Specs {
		typeSpec := typeutil.AsOrEmpty[*ast.TypeSpec](s)

		if typeSpec == nil {
			continue
		}

		_decls, err := g.fromTypeSpec(typeSpec)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		decls = append(decls, _decls...)
	}

	return decls, nil
}

func (g *Generator) fromTypeSpec(typeSpec *ast.TypeSpec) ([]ast.Decl, error) {
	structType := typeutil.AsOrEmpty[*ast.StructType](typeSpec.Type)

	if structType == nil {
		return []ast.Decl{}, nil
	}

	var decls []ast.Decl

	// Generate New function if configured
	if g.config.GenerateNewFunc {
		newFunc := g.newFuncDecl(typeSpec.Name.Name)
		if newFunc != nil {
			decls = append(decls, newFunc)
		}
	}

	// Generate getter and setter methods
	fieldDecls, err := g.fromFieldList(typeSpec.Name.Name, structType.Fields)
	if err != nil {
		return nil, err
	}

	decls = append(decls, fieldDecls...)

	return decls, nil
}

func (g *Generator) fromFieldList(structName string, fieldList *ast.FieldList) ([]ast.Decl, error) {
	var decls []ast.Decl

	for _, f := range fieldList.List {
		_decls, err := g.fromField(structName, f)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		decls = append(decls, _decls...)
	}

	return decls, nil
}

var errInvalidTagValue = errors.New("invalid tag value")

func (g *Generator) fromField(structName string, field *ast.Field) ([]ast.Decl, error) {
	if field.Tag == nil {
		return nil, nil
	}

	tagValue, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	propertyTag := reflect.StructTag(tagValue).Get(g.config.TagName)
	if propertyTag == "" || propertyTag == "-" {
		return []ast.Decl{}, nil
	}

	directives := strings.Split(propertyTag, ",")

	var decls []ast.Decl

	for _, directive := range directives {
		decl, err := g.processDirective(directive, structName, field)
		if err != nil {
			return nil, err
		}

		if decl != nil {
			decls = append(decls, decl)
		}
	}

	return decls, nil
}

func (g *Generator) processDirective(directive, structName string, field *ast.Field) (ast.Decl, error) {
	switch directive {
	case "get":
		return g.getterFuncDecl(structName, field), nil

	case "set":
		return g.setterFuncDecl("Set", structName, field), nil

	case "set=private":
		return g.setterFuncDecl("set", structName, field), nil

	default:
		return nil, errors.Wrapf(errInvalidTagValue, "directive=%s", directive)
	}
}

func (g *Generator) newFuncDecl(structName string) ast.Decl {
	name := astutil.NewIdent("New" + structName)

	funcType := astutil.NewFuncType(
		nil,
		nil,
		astutil.NewFieldList(
			[]*ast.Field{
				astutil.NewField(nil, astutil.NewStarExpr(astutil.NewIdent(structName))),
			},
		),
	)

	body := astutil.NewBlockStmt(
		[]ast.Stmt{
			astutil.NewReturnStmt(
				[]ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X:  astutil.NewCompositeLit(astutil.NewIdent(structName), nil),
					},
				},
			),
		},
	)

	return &ast.FuncDecl{
		Name: name,
		Type: funcType,
		Body: body,
	}
}

func (g *Generator) getterFuncDecl(structName string, field *ast.Field) ast.Decl {
	if len(field.Names) == 0 {
		return nil
	}

	recv := astutil.NewFieldList(
		[]*ast.Field{
			astutil.NewField(
				[]*ast.Ident{
					astutil.NewIdent("t"),
				},
				astutil.NewStarExpr(astutil.NewIdent(structName)),
			),
		},
	)

	name := astutil.NewIdent(
		"Get" + g.prepareFieldName(field.Names[0].Name),
	)

	funcType := astutil.NewFuncType(
		nil,
		nil,
		astutil.NewFieldList(
			[]*ast.Field{
				astutil.NewField(nil, field.Type),
			},
		),
	)

	body := astutil.NewBlockStmt(
		[]ast.Stmt{
			astutil.NewReturnStmt(
				[]ast.Expr{
					astutil.NewSelectorExpr(astutil.NewIdent("t"), astutil.NewIdent(field.Names[0].Name)),
				},
			),
		},
	)

	return &ast.FuncDecl{
		Recv: recv,
		Name: name,
		Type: funcType,
		Body: body,
	}
}

func (g *Generator) setterFuncDecl(verb string, structName string, field *ast.Field) ast.Decl {
	if field.Tag == nil || len(field.Names) == 0 {
		return nil
	}

	tagValue, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		return nil
	}

	validatonTag := reflect.StructTag(tagValue).Get(g.config.ValidationTag)
	if len(validatonTag) > 0 {
		return g.setterFuncWithValidationDecl(verb, structName, field, validatonTag)
	}

	return g.setterFuncNoValidationDecl(verb, structName, field)
}

func (g *Generator) setterFuncNoValidationDecl(verb string, structName string, field *ast.Field) ast.Decl {
	if len(field.Names) == 0 {
		return nil
	}

	recv := astutil.NewFieldList(
		[]*ast.Field{
			astutil.NewField(
				[]*ast.Ident{
					astutil.NewIdent("t"),
				},
				astutil.NewStarExpr(astutil.NewIdent(structName)),
			),
		},
	)

	name := astutil.NewIdent(
		verb + g.prepareFieldName(field.Names[0].Name),
	)

	funcType := g.buildSetterFuncType(field, false)

	body := astutil.NewBlockStmt(
		[]ast.Stmt{
			astutil.NewAssignStmt(
				[]ast.Expr{
					astutil.NewSelectorExpr(astutil.NewIdent("t"), astutil.NewIdent(field.Names[0].Name)),
				},
				token.ASSIGN,
				[]ast.Expr{
					astutil.NewIdent("v"),
				},
			),
		},
	)

	return &ast.FuncDecl{
		Recv: recv,
		Name: name,
		Type: funcType,
		Body: body,
	}
}

func (g *Generator) setterFuncWithValidationDecl(
	verb string, structName string, field *ast.Field, tag string,
) ast.Decl {
	if field.Tag == nil || len(field.Names) == 0 {
		return nil
	}

	tagValue, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		return nil
	}

	validationTag := reflect.StructTag(tagValue).Get(g.config.ValidationTag)
	if len(validationTag) < 1 {
		return nil
	}

	return &ast.FuncDecl{
		Recv: g.buildRecvFieldList(structName),
		Name: astutil.NewIdent(
			verb + g.prepareFieldName(field.Names[0].Name),
		),
		Type: g.buildSetterFuncType(field, true),
		Body: g.buildValidationBody(field, tag),
	}
}

func (g *Generator) buildRecvFieldList(structName string) *ast.FieldList {
	return astutil.NewFieldList(
		[]*ast.Field{
			astutil.NewField(
				[]*ast.Ident{
					astutil.NewIdent("t"),
				},
				astutil.NewStarExpr(astutil.NewIdent(structName)),
			),
		},
	)
}

func (g *Generator) buildSetterFuncType(field *ast.Field, withError bool) *ast.FuncType {
	params := astutil.NewFieldList(
		[]*ast.Field{
			astutil.NewField(
				[]*ast.Ident{
					ast.NewIdent("v"),
				},
				field.Type,
			),
		},
	)

	var results *ast.FieldList

	if withError {
		results = astutil.NewFieldList(
			[]*ast.Field{
				astutil.NewField(nil, astutil.NewIdent("error")),
			},
		)
	}

	return astutil.NewFuncType(nil, params, results)
}

func (g *Generator) buildValidationBody(field *ast.Field, tag string) *ast.BlockStmt {
	callExpr := &ast.CallExpr{
		Fun: astutil.NewIdent(g.config.ValidationFunc),
		Args: []ast.Expr{
			astutil.NewBasicLit(token.STRING, fmt.Sprintf("\"%s\"", field.Names[0].Name)),
			astutil.NewIdent("v"),
			astutil.NewBasicLit(token.STRING, fmt.Sprintf("\"%s\"", tag)),
		},
	}

	return astutil.NewBlockStmt(
		[]ast.Stmt{
			astutil.NewAssignStmt(
				[]ast.Expr{
					astutil.NewIdent("err"),
				},
				token.DEFINE,
				[]ast.Expr{
					callExpr,
				},
			),
			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  astutil.NewIdent("err"),
					Y:  astutil.NewIdent("nil"),
				},
				Body: astutil.NewBlockStmt(
					[]ast.Stmt{
						astutil.NewReturnStmt(
							[]ast.Expr{
								ast.NewIdent("err"),
							},
						),
					},
				),
			},
			astutil.NewAssignStmt(
				[]ast.Expr{
					astutil.NewSelectorExpr(astutil.NewIdent("t"), astutil.NewIdent(field.Names[0].Name)),
				},
				token.ASSIGN,
				[]ast.Expr{
					astutil.NewIdent("v"),
				},
			),
			astutil.NewReturnStmt(
				[]ast.Expr{
					astutil.NewIdent("nil"),
				},
			),
		},
	)
}

var camelHeadPattern = regexp.MustCompile(`^[a-z]+`)

func (g *Generator) prepareFieldName(name string) string {
	head := camelHeadPattern.FindString(name)

	if len(head) > 0 {
		head = cases.Title(language.Und).String(head)

		for _, s := range g.config.Initialism {
			if cases.Title(language.Und).String(s) == head {
				head = strings.ToUpper(head)
				break
			}
		}

		name = camelHeadPattern.ReplaceAllString(name, head)
	}

	return name
}
