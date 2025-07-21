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
	PropertyTag    string
	Initialism     []string
	ValidationFunc string
	ValidationTag  string
	NewFunc        bool
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
	if g.config.NewFunc {
		newFunc := g.newFuncDecl(typeSpec.Name.Name, structType.Fields)
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

	propertyTag := reflect.StructTag(tagValue).Get(g.config.PropertyTag)
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

func (g *Generator) newFuncDecl(structName string, fieldList *ast.FieldList) ast.Decl {
	params, assignments, hasValidation := g.buildNewFuncParams(fieldList)
	funcType := g.buildNewFuncType(structName, params, hasValidation)
	body := g.buildNewFuncBody(structName, assignments, hasValidation)

	return &ast.FuncDecl{
		Name: astutil.NewIdent("New" + structName),
		Type: funcType,
		Body: body,
	}
}

func (g *Generator) buildNewFuncParams(fieldList *ast.FieldList) ([]*ast.Field, []ast.Stmt, bool) {
	fieldCount := len(fieldList.List)

	const estimatedStmtsPerField = 2

	constructorParams := make([]*ast.Field, 0, fieldCount)
	// Estimate estimatedStmtsPerField statements per field
	assignments := make([]ast.Stmt, 0, fieldCount*estimatedStmtsPerField)

	var (
		hasValidation bool
		errDeclared   bool
	)

	for _, field := range fieldList.List {
		if len(field.Names) == 0 {
			continue
		}

		param, assignment, validation := g.processFieldForConstructor(field, &errDeclared)

		constructorParams = append(constructorParams, param)
		if assignment != nil {
			assignments = append(assignments, assignment...)
		}

		if validation {
			hasValidation = true
		}
	}

	return constructorParams, assignments, hasValidation
}

func (g *Generator) processFieldForConstructor(field *ast.Field, errDeclared *bool) (*ast.Field, []ast.Stmt, bool) {
	paramName := field.Names[0].Name
	param := astutil.NewField(
		[]*ast.Ident{astutil.NewIdent(paramName)},
		field.Type,
	)

	if field.Tag == nil {
		return param, g.createDirectAssignment(paramName), false
	}

	tagValue, err := strconv.Unquote(field.Tag.Value)
	if err != nil {
		return param, nil, false
	}

	propertyTag := reflect.StructTag(tagValue).Get(g.config.PropertyTag)
	if propertyTag == "" || propertyTag == "-" {
		return param, g.createDirectAssignment(paramName), false
	}

	_, assignment, validation := g.processFieldForNewFunc(field, tagValue, propertyTag)

	const minStatementsForValidation = 2
	if assignment != nil && validation && len(assignment) >= minStatementsForValidation {
		g.adjustErrDeclaration(assignment, errDeclared)
	}

	return param, assignment, validation
}

func (g *Generator) createDirectAssignment(paramName string) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				astutil.NewSelectorExpr(astutil.NewIdent("s"), astutil.NewIdent(paramName)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{astutil.NewIdent(paramName)},
		},
	}
}

func (g *Generator) adjustErrDeclaration(assignment []ast.Stmt, errDeclared *bool) {
	if *errDeclared {
		// Change := to = for subsequent err assignments
		if assignStmt, ok := assignment[0].(*ast.AssignStmt); ok {
			assignStmt.Tok = token.ASSIGN
		}
	} else {
		*errDeclared = true
	}
}

func (g *Generator) processFieldForNewFunc(
	field *ast.Field, tagValue, propertyTag string,
) (*ast.Field, []ast.Stmt, bool) {
	paramName := field.Names[0].Name
	param := astutil.NewField(
		[]*ast.Ident{astutil.NewIdent(paramName)},
		field.Type,
	)

	// If no property tag, use direct assignment
	if propertyTag == "" {
		return param, g.createDirectAssignment(paramName), false
	}

	directives := strings.Split(propertyTag, ",")
	hasSetter, isPrivateSetter := g.parseSetterDirectives(directives)

	// If there's no setter, use direct assignment
	if !hasSetter {
		return param, g.createDirectAssignment(paramName), false
	}

	return g.buildSetterAssignment(param, paramName, field, tagValue, isPrivateSetter)
}

func (g *Generator) parseSetterDirectives(directives []string) (bool, bool) {
	var hasSetter, isPrivateSetter bool

	for _, directive := range directives {
		switch directive {
		case "set":
			hasSetter = true
		case "set=private":
			hasSetter = true
			isPrivateSetter = true
		}
	}

	return hasSetter, isPrivateSetter
}

func (g *Generator) buildSetterAssignment(
	param *ast.Field, paramName string, field *ast.Field, tagValue string, isPrivateSetter bool,
) (*ast.Field, []ast.Stmt, bool) {
	validationTag := reflect.StructTag(tagValue).Get(g.config.ValidationTag)

	methodName := "Set" + g.prepareFieldName(field.Names[0].Name)
	if isPrivateSetter {
		methodName = "set" + g.prepareFieldName(field.Names[0].Name)
	}

	if len(validationTag) > 0 {
		// Add validation call
		return param, []ast.Stmt{
			astutil.NewAssignStmt(
				[]ast.Expr{astutil.NewIdent("err")},
				token.DEFINE,
				[]ast.Expr{
					&ast.CallExpr{
						Fun:  astutil.NewSelectorExpr(astutil.NewIdent("s"), astutil.NewIdent(methodName)),
						Args: []ast.Expr{astutil.NewIdent(paramName)},
					},
				},
			),
			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					Op: token.NEQ,
					X:  astutil.NewIdent("err"),
					Y:  astutil.NewIdent("nil"),
				},
				Body: astutil.NewBlockStmt([]ast.Stmt{
					astutil.NewReturnStmt([]ast.Expr{
						astutil.NewIdent("nil"),
						astutil.NewIdent("err"),
					}),
				}),
			},
		}, true
	}

	// Direct setter call
	return param, []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  astutil.NewSelectorExpr(astutil.NewIdent("s"), astutil.NewIdent(methodName)),
				Args: []ast.Expr{astutil.NewIdent(paramName)},
			},
		},
	}, false
}

func (g *Generator) buildNewFuncType(structName string, params []*ast.Field, hasValidation bool) *ast.FuncType {
	if hasValidation {
		return astutil.NewFuncType(
			nil,
			astutil.NewFieldList(params),
			astutil.NewFieldList([]*ast.Field{
				astutil.NewField(nil, astutil.NewStarExpr(astutil.NewIdent(structName))),
				astutil.NewField(nil, astutil.NewIdent("error")),
			}),
		)
	}

	return astutil.NewFuncType(
		nil,
		astutil.NewFieldList(params),
		astutil.NewFieldList([]*ast.Field{
			astutil.NewField(nil, astutil.NewStarExpr(astutil.NewIdent(structName))),
		}),
	)
}

func (g *Generator) buildNewFuncBody(structName string, assignments []ast.Stmt, hasValidation bool) *ast.BlockStmt {
	// Special case for empty struct
	if len(assignments) == 0 {
		return astutil.NewBlockStmt([]ast.Stmt{
			astutil.NewReturnStmt([]ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X:  astutil.NewCompositeLit(astutil.NewIdent(structName), nil),
				},
			}),
		})
	}

	var bodyStmts []ast.Stmt

	// Create struct instance
	bodyStmts = append(bodyStmts,
		astutil.NewAssignStmt(
			[]ast.Expr{astutil.NewIdent("s")},
			token.DEFINE,
			[]ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X:  astutil.NewCompositeLit(astutil.NewIdent(structName), nil),
				},
			},
		),
	)

	// Add field assignments
	bodyStmts = append(bodyStmts, assignments...)

	// Add return statement
	if hasValidation {
		bodyStmts = append(bodyStmts,
			astutil.NewReturnStmt([]ast.Expr{
				astutil.NewIdent("s"),
				astutil.NewIdent("nil"),
			}),
		)
	} else {
		bodyStmts = append(bodyStmts,
			astutil.NewReturnStmt([]ast.Expr{
				astutil.NewIdent("s"),
			}),
		)
	}

	return astutil.NewBlockStmt(bodyStmts)
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
