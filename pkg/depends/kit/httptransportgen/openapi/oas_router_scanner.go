package openapi

import (
	"context"
	"go/ast"
	"go/types"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func NewRouterScanner(pkg *pkgx.Pkg) *RouterScanner {
	rs := &RouterScanner{
		pkg:     pkg,
		routers: map[*types.Var]*Router{},
		os:      NewOperatorScanner(pkg),
	}
	rs.init()
	return rs
}

type RouterScanner struct {
	pkg     *pkgx.Pkg
	routers map[*types.Var]*Router
	os      *OperatorScanner
}

func (rs *RouterScanner) init() {
	for _, pkg := range rs.pkg.Imports() {
		for id, obj := range pkg.TypesInfo.Defs {
			v, ok := obj.(*types.Var)
			if !ok {
				continue
			}
			if v == nil || strings.HasSuffix(v.Pkg().Path(), PkgPathKit) ||
				!isRouterType(v.Type()) {
				continue
			}
			r := NewRouter(v)
			ast.Inspect(
				id.Obj.Decl.(ast.Node),
				func(node ast.Node) bool {
					switch call := node.(type) {
					case *ast.CallExpr:
						r.AppendOperators(
							rs.OperatorTypeNamesFromArgs(
								pkgx.New(pkg),
								call.Args...,
							)...,
						)
						return false
					}
					return true
				},
			)

			rs.routers[v] = r
		}
	}
	for _, pkg := range rs.pkg.Imports() {
		for sel, selection := range pkg.TypesInfo.Selections {
			if selection.Obj() == nil {
				continue
			}
			f, ok := selection.Obj().(*types.Func)
			if !ok {
				continue
			}
			recv := f.Type().(*types.Signature).Recv()
			if recv == nil || !isRouterType(recv.Type()) {
				continue
			}
			for tv, r := range rs.routers {
				if sel.Sel.Name != "Register" {
					continue
				}
				if tv != pkg.TypesInfo.ObjectOf(pkgx.GetIdentChainOfCallFunc(sel)[0]) {
					continue
				}
				file := rs.pkg.FileOf(sel)
				ast.Inspect(file, func(node ast.Node) bool {
					call, ok := node.(*ast.CallExpr)
					if !ok {
						return true
					}
					if call.Fun == sel {
						ident := call.Args[0]
						switch v := ident.(type) {
						case *ast.Ident:
							arg := pkg.TypesInfo.ObjectOf(v).(*types.Var)
							if router, ok := rs.routers[arg]; ok {
								r.Register(router)
							}
						case *ast.SelectorExpr:
							arg := pkg.TypesInfo.ObjectOf(v.Sel).(*types.Var)
							if router, ok := rs.routers[arg]; ok {
								r.Register(router)
							}
						case *ast.CallExpr:
							r.With(rs.OperatorTypeNamesFromArgs(pkgx.New(pkg), v.Args...)...)
						}
						return false
					}
					return true

				})
			}
		}
	}
}

func (rs *RouterScanner) Router(tn *types.Var) *Router {
	return rs.routers[tn]
}

func (rs *RouterScanner) OperatorTypeNamesFromArgs(pkg *pkgx.Pkg, args ...ast.Expr) []*OperatorWithTypeName {
	ops := make([]*OperatorWithTypeName, 0)

	for _, arg := range args {
		op := rs.OperatorTypeNameFromType(pkg.TypesInfo.TypeOf(arg))

		if op == nil {
			continue
		}

		// modify meta if httptransport.Group() or httptransport.BasePath()
		if call, ok := arg.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if isFromHttpTransport(pkg.TypesInfo.ObjectOf(sel.Sel).Type()) {
					switch sel.Sel.Name {
					case "BasePath":
						switch v := call.Args[0].(type) {
						case *ast.BasicLit:
							op.BasePath, _ = strconv.Unquote(v.Value)
						}
					case "Group":
						switch v := call.Args[0].(type) {
						case *ast.BasicLit:
							op.Path, _ = strconv.Unquote(v.Value)
						}
					}
				}
			}
		}

		// handle interface WithMiddleOperators
		if op.TypeName != nil {
			_, res, ok := AssertIfByMtdNameAndResCntInPkg(
				op.TypeName.Type(), rs.pkg, "MiddleOperators", 1,
			)
			if ok {
				for _, v := range res[0] {
					if lit, ok := v.Expr.(*ast.CompositeLit); ok {
						ops = append(
							ops,
							rs.OperatorTypeNamesFromArgs(pkg, lit.Elts...)...,
						)
					}

				}

			}
		}

		ops = append(ops, op)
	}

	return ops
}

func (rs *RouterScanner) OperatorTypeNameFromType(typ types.Type) *OperatorWithTypeName {
	switch t := typ.(type) {
	case *types.Pointer:
		return rs.OperatorTypeNameFromType(t.Elem())
	case *types.Named:
		tn := t.Obj()
		if op := rs.os.Operator(context.Background(), tn); op != nil {
			return &OperatorWithTypeName{
				Operator: op,
				TypeName: tn,
			}
		}

		return nil
	default:
		return nil
	}
}
