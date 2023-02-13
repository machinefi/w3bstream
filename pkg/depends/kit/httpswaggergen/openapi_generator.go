package httpswaggergen

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/types"
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
		pkg:           pkg,
		openapi:       oas.NewOpenAPI(),
		routerScanner: NewRouterScanner(pkg),
	}
}

type OpenAPIGenerator struct {
	pkg           *pkgx.Pkg
	openapi       *oas.OpenAPI
	routerScanner *RouterScanner
}

func rootRouter(pkgInfo *pkgx.Pkg, callExpr *ast.CallExpr) *types.Var {
	if len(callExpr.Args) > 0 {
		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if typesFunc, ok := pkgInfo.TypesInfo.ObjectOf(selectorExpr.Sel).(*types.Func); ok {
				if signature, ok := typesFunc.Type().(*types.Signature); ok {
					if isRouterType(signature.Params().At(0).Type()) {
						if selectorExpr.Sel.Name == "Run" || selectorExpr.Sel.Name == "Serve" {
							switch node := callExpr.Args[0].(type) {
							case *ast.SelectorExpr:
								return pkgInfo.TypesInfo.ObjectOf(node.Sel).(*types.Var)
							case *ast.Ident:
								return pkgInfo.TypesInfo.ObjectOf(node).(*types.Var)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (g *OpenAPIGenerator) Scan(ctx context.Context) {
	defer func() {
		g.routerScanner.operatorScanner.BindSchemas(g.openapi)
	}()

	for ident, def := range g.pkg.TypesInfo.Defs {
		if typFunc, ok := def.(*types.Func); ok {
			if typFunc.Name() != "main" {
				continue
			}

			ast.Inspect(ident.Obj.Decl.(*ast.FuncDecl), func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.CallExpr:
					if rootRouterVar := rootRouter(g.pkg, n); rootRouterVar != nil {
						router := g.routerScanner.Router(rootRouterVar)

						routes := router.Routes()

						operationIDs := map[string]*Route{}

						for _, route := range routes {
							method := route.Method()

							operation := g.OperationByOperatorTypes(method, route.Operators...)

							if _, exists := operationIDs[operation.OperationId]; exists {
								panic(errors.Errorf("operationID %s should be unique", operation.OperationId))
							}

							operationIDs[operation.OperationId] = route

							g.openapi.AddOperation(oas.HttpMethod(strings.ToLower(method)), g.patchPath(route.Path(), operation), operation)
						}
					}
				}
				return true
			})
			return
		}
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
	data, err := json.MarshalIndent(g.openapi, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(file, data, os.ModePerm)
	log.Printf("generated openapi spec into %s", color.MagentaString(file))
}
