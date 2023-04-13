package openapi

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/pkg/errors"
)

func NewOpenAPIGenerator(pkg *pkgx.Pkg) *OpenAPIGenerator {
	return &OpenAPIGenerator{
		pkg: pkg,
		oas: oas.NewOpenAPI(),
		rs:  NewRouterScanner(pkg),
	}
}

type OpenAPIGenerator struct {
	pkg *pkgx.Pkg
	oas *oas.OpenAPI
	rs  *RouterScanner
}

func rootRouter(pkg *pkgx.Pkg, call *ast.CallExpr) *types.Var {
	if len(call.Args) == 0 {
		return nil
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	tf, ok := pkg.TypesInfo.ObjectOf(sel.Sel).(*types.Func)
	if !ok {
		return nil
	}
	sig, ok := tf.Type().(*types.Signature)
	if !ok {
		return nil
	}
	if !isRouterType(sig.Params().At(0).Type()) {
		return nil
	}
	if sel.Sel.Name == "Run" || sel.Sel.Name == "Serve" {
		switch node := call.Args[0].(type) {
		case *ast.SelectorExpr:
			return pkg.TypesInfo.ObjectOf(node.Sel).(*types.Var)
		case *ast.Ident:
			return pkg.TypesInfo.ObjectOf(node).(*types.Var)
		}
	}
	return nil
}

func (g *OpenAPIGenerator) Scan(ctx context.Context) {
	defer func() {
		g.rs.os.BindSchemas(g.oas)
	}()

	for ident, def := range g.pkg.TypesInfo.Defs {
		tf, ok := def.(*types.Func)
		if !ok || tf.Name() != "main" {
			continue
		}
		ast.Inspect(ident.Obj.Decl.(*ast.FuncDecl), func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.CallExpr:
				if rootRouterVar := rootRouter(g.pkg, n); rootRouterVar != nil {
					router := g.rs.Router(rootRouterVar)

					routes := router.Routes()

					operationIDs := map[string]*Route{}

					for _, route := range routes {
						method := route.Method()

						operation := g.OperationByOperatorTypes(method, route.Operators...)

						if _, exists := operationIDs[operation.OperationId]; exists {
							panic(errors.Errorf("operationID %s should be unique", operation.OperationId))
						}

						operationIDs[operation.OperationId] = route

						g.oas.AddOperation(oas.HttpMethod(strings.ToLower(method)), g.patchPath(route.Path(), operation), operation)
					}
				}
			}
			return true
		})
		return

	}
}

var reHttpRouterPath = regexp.MustCompile("/:([^/]+)")

func (g *OpenAPIGenerator) patchPath(openapiPath string, operation *oas.Operation) string {
	return reHttpRouterPath.ReplaceAllStringFunc(openapiPath, func(str string) string {
		name := reHttpRouterPath.FindAllStringSubmatch(str, -1)[0][1]

		var isParameterDefined = false

		for _, parameter := range operation.Parameters {
			if parameter.In == "path" && parameter.Name == name {
				isParameterDefined = true
			}
		}

		if isParameterDefined {
			return "/{" + name + "}"
		}

		return "/0"
	})
}

func (g *OpenAPIGenerator) OperationByOperatorTypes(method string, operatorTypes ...*OperatorWithTypeName) *oas.Operation {
	operation := &oas.Operation{}

	length := len(operatorTypes)

	for idx := range operatorTypes {
		operatorTypes[idx].BindOperation(method, operation, idx == length-1)
	}

	return operation
}

func (g *OpenAPIGenerator) Output(cwd string) {
	file := filepath.Join(cwd, "openapi.json")
	data, err := json.MarshalIndent(g.oas, "", "  ")
	if err != nil {
		return
	}
	_ = ioutil.WriteFile(file, data, os.ModePerm)
	log.Printf("generated oas spec into %s", color.MagentaString(file))
}
