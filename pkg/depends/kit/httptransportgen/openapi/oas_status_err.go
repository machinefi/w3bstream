package openapi

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strconv"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusxgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func NewStatusErrScanner(pkg *pkgx.Pkg) *StatusErrScanner {
	s := &StatusErrScanner{
		pkg:  pkg,
		defs: map[*types.Named][]*statusx.StatusErr{},
		used: map[*types.Func][]*statusx.StatusErr{},
	}
	s.init()
	return s
}

type StatusErrScanner struct {
	Type *types.Named
	pkg  *pkgx.Pkg
	defs map[*types.Named][]*statusx.StatusErr
	used map[*types.Func][]*statusx.StatusErr
}

func (s *StatusErrScanner) init() {
	pkg := s.pkg.PkgByPath(PkgPathStatusErr)
	if pkg == nil {
		return
	}
	s.Type = pkgx.New(pkg).TypeName("StatusErr").Type().(*types.Named)

	scanner := statusxgen.NewScanner(s.pkg)

	for _, pkgInfo := range s.pkg.Imports() {
		for _, obj := range pkgInfo.TypesInfo.Defs {
			if tn, ok := obj.(*types.TypeName); ok {
				if isStatusError(tn.Type(), pkg) {
					s.defs[tn.Type().(*types.Named)] = scanner.StatusError(tn)
				}
			}
		}
	}
}

func (s *StatusErrScanner) StatusErrorsInFunc(tf *types.Func) []*statusx.StatusErr {
	if tf == nil {
		return nil
	}

	if errs, ok := s.used[tf]; ok {
		return errs
	}

	s.used[tf] = []*statusx.StatusErr{}

	pkg := pkgx.New(s.pkg.PkgByPath(tf.Pkg().Path()))
	fd := pkg.FuncDeclOf(tf)

	if fd == nil {
		return s.used[tf]
	}

	ast.Inspect(fd, func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.CallExpr:
			idents := pkgx.GetIdentChainOfCallFunc(v.Fun)
			if len(idents) > 0 {
				call := idents[len(idents)-1]
				obj := pkg.TypesInfo.ObjectOf(call)
				if obj != nil {
					// pick status errors from statusx.Wrap
					if call.Name == "Wrap" && obj.Pkg().Path() == PkgPathStatusErr {
						code := 0
						key := ""
						msg := ""
						desc := make([]string, 0)

						for i, arg := range v.Args[1:] {
							tv, err := pkg.Eval(arg)
							if err != nil {
								continue
							}
							switch i {
							case 0: // code
								code, _ = strconv.Atoi(tv.Value.String())
							case 1: // key
								key, _ = strconv.Unquote(tv.Value.String())
							case 2: // msg
								msg, _ = strconv.Unquote(tv.Value.String())
							default:
								d, _ := strconv.Unquote(tv.Value.String())
								desc = append(desc, d)
							}
						}
						if code > 0 {
							if msg == "" {
								msg = key
							}
							s.append(tf,
								statusx.Wrap(
									errors.New(""), code, key,
									append([]string{msg}, desc...)...,
								),
							)
						}

					}
				}
				// Deprecated old code defined
				if obj != nil && obj.Pkg() != nil &&
					obj.Pkg().Path() == s.Type.Obj().Pkg().Path() {
					for i := range idents {
						s.mayAddByObject(tf, pkg.TypesInfo.ObjectOf(idents[i]))
					}
					return false
				}

				if next, ok := obj.(*types.Func); ok && next != tf && next.Pkg() != nil {
					s.append(tf, s.StatusErrorsInFunc(next)...)
				}
			}
		case *ast.Ident:
			s.mayAddByObject(tf, pkg.TypesInfo.ObjectOf(v))
		}
		return true
	})
	doc := pkgx.StringifyCommentGroup(fd.Doc)
	s.append(tf, PickStatusErrorsFromDoc(doc)...)
	return s.used[tf]
}

func (s *StatusErrScanner) mayAddByObject(tf *types.Func, obj types.Object) {
	if obj == nil {
		return
	}
	if tc, ok := obj.(*types.Const); ok {
		if named, ok := tc.Type().(*types.Named); ok {
			if errs, ok := s.defs[named]; ok {
				for i := range errs {
					if errs[i].Key == tc.Name() {
						s.append(tf, errs[i])
					}
				}
			}
		}
	}
}

func (s *StatusErrScanner) append(tf *types.Func, errs ...*statusx.StatusErr) {
	m := map[string]*statusx.StatusErr{}

	for _, se := range append(s.used[tf], errs...) {
		m[fmt.Sprintf("%s%d", se.Key, se.Code)] = se
	}

	next := make([]*statusx.StatusErr, 0)
	for k := range m {
		next = append(next, m[k])
	}

	sort.Slice(next, func(i, j int) bool {
		return next[i].Code < next[j].Code
	})

	s.used[tf] = next
}
