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

type GeneratorConfig struct {
	TagName        string
	Initialism     []string
	Validate       bool
	ValidationFunc string
	ValidationTag  string
}

type Generator struct {
	config *GeneratorConfig
}

func NewGenerator(config *GeneratorConfig) *Generator {
	return &Generator{
		config: config,
	}
}

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

	return g.fromFieldList(typeSpec.Name.Name, structType.Fields)
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

func (g *Generator) fromField(structName string, field *ast.Field) ([]ast.Decl, error) {
	if field.Tag == nil {
		return nil, nil
	}

	tagValue, _ := strconv.Unquote(field.Tag.Value)
	propertyTag := reflect.StructTag(tagValue).Get(g.config.TagName)

	directives := strings.Split(propertyTag, ",")
	if len(directives) == 0 || (len(directives) == 1 && (directives[0] == "" || directives[0] == "-")) {
		return []ast.Decl{}, nil
	}

	var decls []ast.Decl

	for _, directive := range directives {
		switch directive {
		case "get":
			decl := g.getterFuncDecl(structName, field)
			if decl != nil {
				decls = append(decls, decl)
			}

		case "set":
			decl := g.setterFuncDecl("Set", structName, field)
			if decl != nil {
				decls = append(decls, decl)
			}

		case "set=private":
			decl := g.setterFuncDecl("set", structName, field)
			if decl != nil {
				decls = append(decls, decl)
			}

		default:
			return nil, fmt.Errorf("invalid tag value '%s'", directive)
		}
	}

	return decls, nil
}

func (g *Generator) getterFuncDecl(structName string, field *ast.Field) ast.Decl {
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
	if field.Tag == nil {
		return nil
	}

	tagValue, _ := strconv.Unquote(field.Tag.Value)
	validatonTag := reflect.StructTag(tagValue).Get(g.config.ValidationTag)

	if len(validatonTag) > 0 {
		return g.setterFuncWithValidationDecl(verb, structName, field, validatonTag)
	}

	return g.setterFuncNoValidationDecl(verb, structName, field)
}

func (g *Generator) setterFuncNoValidationDecl(verb string, structName string, field *ast.Field) ast.Decl {
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
	funcType := astutil.NewFuncType(
		nil,
		astutil.NewFieldList(
			[]*ast.Field{
				astutil.NewField(
					[]*ast.Ident{
						ast.NewIdent("v"),
					},
					field.Type,
				),
			},
		),
		nil,
	)
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

func (g *Generator) setterFuncWithValidationDecl(verb string, structName string, field *ast.Field, tag string) ast.Decl {
	if field.Tag == nil {
		return nil
	}

	t1, _ := strconv.Unquote(field.Tag.Value)
	t2 := reflect.StructTag(t1).Get(g.config.ValidationTag)

	if len(t2) < 1 {
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
	funcType := astutil.NewFuncType(
		nil,
		astutil.NewFieldList(
			[]*ast.Field{
				astutil.NewField(
					[]*ast.Ident{
						ast.NewIdent("v"),
					},
					field.Type,
				),
			},
		),
		astutil.NewFieldList(
			[]*ast.Field{
				astutil.NewField(nil, astutil.NewIdent("error")),
			},
		),
	)
	body := astutil.NewBlockStmt(
		[]ast.Stmt{
			astutil.NewAssignStmt(
				[]ast.Expr{
					astutil.NewIdent("err"),
				},
				token.DEFINE,
				[]ast.Expr{
					astutil.NewIdent(fmt.Sprintf("%s(v, \"%s\")", g.config.ValidationFunc, tag)),
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

	return &ast.FuncDecl{
		Recv: recv,
		Name: name,
		Type: funcType,
		Body: body,
	}
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
