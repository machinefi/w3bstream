package statusxgen

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func NewScanner(pkg *pkgx.Pkg) *Scanner {
	return &Scanner{
		pkg: pkg,
	}
}

type Scanner struct {
	pkg          *pkgx.Pkg
	StatusErrors map[*types.TypeName][]*statusx.StatusErr
}

func sortedStatusErrList(list []*statusx.StatusErr) []*statusx.StatusErr {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Code < list[j].Code
	})
	return list
}

func (s *Scanner) StatusError(tn *types.TypeName) []*statusx.StatusErr {
	if tn == nil {
		return nil
	}

	if es, ok := s.StatusErrors[tn]; ok {
		sort.Slice(es, func(i, j int) bool {
			return es[i].Code < es[j].Code
		})
		return es
	}

	if !strings.Contains(tn.Type().Underlying().String(), "int") {
		panic(fmt.Errorf("status error type underlying must be an int or uint, but got %s", tn.String()))
	}

	pkg := s.pkg.PkgByPath(tn.Pkg().Path())
	if pkg == nil {
		return nil
	}

	serviceCode := 0

	method, ok := typesx.FromGoType(tn.Type()).MethodByName("ServiceCode")
	if ok {
		results, n := s.pkg.FuncResultsOf(method.(*typesx.GoMethod).Func)
		if n == 1 {
			ret := results[0][0]
			if ret.IsValue() {
				if i, err := strconv.ParseInt(ret.Value.String(), 10, 64); err == nil {
					serviceCode = int(i)
				}
			}
		}
	}

	for ident, def := range pkg.TypesInfo.Defs {
		typeConst, ok := def.(*types.Const)
		if !ok {
			continue
		}
		if typeConst.Type() != tn.Type() {
			continue
		}

		key := typeConst.Name()
		code, _ := strconv.ParseInt(typeConst.Val().String(), 10, 64)

		msg, canBeTalkError := ParseStatusErrMsg(ident.Obj.Decl.(*ast.ValueSpec).Doc.Text())

		s.add(tn, key, msg, int(code)+serviceCode, canBeTalkError)
	}

	lst := s.StatusErrors[tn]
	sort.Slice(lst, func(i, j int) bool {
		return lst[i].Code < lst[j].Code
	})
	return lst
}

func ParseStatusErrMsg(s string) (string, bool) {
	firstLine := strings.Split(strings.TrimSpace(s), "\n")[0]

	prefix := "@errTalk "
	if strings.HasPrefix(firstLine, prefix) {
		return firstLine[len(prefix):], true
	}
	return firstLine, false
}

func (s *Scanner) add(tn *types.TypeName, key, msg string, code int, canBeTalk bool) {
	if s.StatusErrors == nil {
		s.StatusErrors = map[*types.TypeName][]*statusx.StatusErr{}
	}

	se := statusx.NewStatusErr(key, code, msg)
	if canBeTalk {
		se = se.EnableErrTalk()
	}
	s.StatusErrors[tn] = append(s.StatusErrors[tn], se)
}
