package openapi

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func NewGenerator(pkg *pkgx.Pkg) *Generator {
	return &Generator{
		pkg: pkg,
		oas: oas.NewOpenAPI(),
		rs:  NewRouterScanner(pkg),
	}
}

type Generator struct {
	pkg *pkgx.Pkg
	oas *oas.OpenAPI
	rs  *RouterScanner
}

func root(pkg *pkgx.Pkg, call *ast.CallExpr) *types.Var {
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

func (g *Generator) Scan(ctx context.Context) {
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
				if rv := root(g.pkg, n); rv != nil {
					router := g.rs.Router(rv)
					routes := router.Routes()
					ops := map[string]*Route{}

					for _, route := range routes {
						mtd := route.Method()
						op := g.OperationByOperatorTypes(mtd, route.Operators...)

						if _, exists := ops[op.OperationId]; exists {
							panic(errors.Errorf("operationID %s should be unique", op.OperationId))
						}

						ops[op.OperationId] = route
						g.oas.AddOperation(
							oas.HttpMethod(strings.ToLower(mtd)),
							PatchRouterPath(route.Path(), op), op,
						)
					}
				}
			}
			return true
		})
		return
	}
}

func (g *Generator) OperationByOperatorTypes(mtd string, ops ...*OperatorWithTypeName) *oas.Operation {
	op := &oas.Operation{}

	length := len(ops)

	for idx := range ops {
		ops[idx].BindOperation(mtd, op, idx == length-1)
	}

	return op
}

func (g *Generator) Output(cwd string) {
	file := filepath.Join(cwd, "openapi.json")
	data, err := json.MarshalIndent(g.oas, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(file, data, os.ModePerm)
	log.Printf("generated oas spec into %s", color.MagentaString(file))
}
