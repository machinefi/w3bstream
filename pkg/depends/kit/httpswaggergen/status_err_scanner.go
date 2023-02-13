package httpswaggergen

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusxgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/pkg/errors"
)

func NewStatusErrScanner(pkg *pkgx.Pkg) *StatusErrScanner {
	statusErrorScanner := &StatusErrScanner{
		pkg:              pkg,
		statusErrorTypes: map[*types.Named][]*statusx.StatusErr{},
		errorsUsed:       map[*types.Func][]*statusx.StatusErr{},
	}

	statusErrorScanner.init()

	return statusErrorScanner
}

type StatusErrScanner struct {
	StatusErrType    *types.Named
	pkg              *pkgx.Pkg
	statusErrorTypes map[*types.Named][]*statusx.StatusErr
	errorsUsed       map[*types.Func][]*statusx.StatusErr
}

var statusErrPkgPath = reflect.TypeOf(statusx.StatusErr{}).PkgPath()

func (es *StatusErrScanner) StatusErrorsInFunc(typeFunc *types.Func) []*statusx.StatusErr {
	if typeFunc == nil {
		return nil
	}

	if statusErrList, ok := es.errorsUsed[typeFunc]; ok {
		return statusErrList
	}

	es.errorsUsed[typeFunc] = []*statusx.StatusErr{}

	pkg := pkgx.New(es.pkg.PkgByPath(typeFunc.Pkg().Path()))

	funcDecl := pkg.FuncDeclOf(typeFunc)

	if funcDecl != nil {
		ast.Inspect(funcDecl, func(node ast.Node) bool {
			switch v := node.(type) {
			case *ast.CallExpr:
				identList := pkgx.GetIdentChainOfCallFunc(v.Fun)
				if len(identList) > 0 {
					callIdent := identList[len(identList)-1]
					obj := pkg.TypesInfo.ObjectOf(callIdent)

					if obj != nil {
						// pick status errors from statusx.Wrap
						if callIdent.Name == "Wrap" && obj.Pkg().Path() == statusErrPkgPath {

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

								es.appendStateErrs(typeFunc, statusx.Wrap(errors.New(""), code, key, append([]string{msg}, desc...)...))
							}

						}
					}

					// Deprecated old code defined
					if obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == es.StatusErrType.Obj().Pkg().Path() {
						for i := range identList {
							es.mayAddStateErrorByObject(typeFunc, pkg.TypesInfo.ObjectOf(identList[i]))
						}
						return false
					}

					if nextTypeFunc, ok := obj.(*types.Func); ok && nextTypeFunc != typeFunc && nextTypeFunc.Pkg() != nil {
						es.appendStateErrs(typeFunc, es.StatusErrorsInFunc(nextTypeFunc)...)
					}
				}
			case *ast.Ident:
				es.mayAddStateErrorByObject(typeFunc, pkg.TypesInfo.ObjectOf(v))
			}
			return true
		})

		doc := pkgx.StringifyCommentGroup(funcDecl.Doc)
		es.appendStateErrs(typeFunc, pickStatusErrorsFromDoc(doc)...)
	}

	return es.errorsUsed[typeFunc]
}

func (es *StatusErrScanner) mayAddStateErrorByObject(typeFunc *types.Func, obj types.Object) {
	if obj == nil {
		return
	}
	if typeConst, ok := obj.(*types.Const); ok {
		if named, ok := typeConst.Type().(*types.Named); ok {
			if errs, ok := es.statusErrorTypes[named]; ok {
				for i := range errs {
					if errs[i].Key == typeConst.Name() {
						es.appendStateErrs(typeFunc, errs[i])
					}
				}
			}
		}
	}
}

func (es *StatusErrScanner) appendStateErrs(typeFunc *types.Func, statusErrs ...*statusx.StatusErr) {
	m := map[string]*statusx.StatusErr{}

	errs := append(es.errorsUsed[typeFunc], statusErrs...)
	for i := range errs {
		s := errs[i]
		m[fmt.Sprintf("%s%d", s.Key, s.Code)] = s
	}

	next := make([]*statusx.StatusErr, 0)
	for k := range m {
		next = append(next, m[k])
	}

	sort.Slice(next, func(i, j int) bool {
		return next[i].Code < next[j].Code
	})

	es.errorsUsed[typeFunc] = next
}

func (es *StatusErrScanner) init() {
	pkg := es.pkg.PkgByPath(pkgPathStatusx)
	if pkg == nil {
		return
	}

	es.StatusErrType = pkgx.New(pkg).TypeName("StatusErr").Type().(*types.Named)
	goTypeStatusError := pkgx.New(pkg).TypeName("Error").Type().Underlying().(*types.Interface)

	isStatusError := func(typ *types.TypeName) bool {
		return types.Implements(typ.Type(), goTypeStatusError)
	}

	s := statusxgen.NewScanner(es.pkg)

	for _, pkgInfo := range es.pkg.Imports() {
		for _, obj := range pkgInfo.TypesInfo.Defs {
			if typName, ok := obj.(*types.TypeName); ok {
				if isStatusError(typName) {
					es.statusErrorTypes[typName.Type().(*types.Named)] = s.StatusError(typName)
				}
			}
		}
	}
}

func pickStatusErrorsFromDoc(doc string) []*statusx.StatusErr {
	statusErrorList := make([]*statusx.StatusErr, 0)

	lines := strings.Split(doc, "\n")

	for _, line := range lines {
		if line != "" {
			if statusErr, err := statusx.ParseStatusErrSummary(line); err == nil {
				statusErrorList = append(statusErrorList, statusErr)
			}
		}
	}

	return statusErrorList
}
