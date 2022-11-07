package enumgen

import (
	"go/ast"
	"go/constant"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

type Scanner struct {
	pkg *pkgx.Pkg
	res map[*types.TypeName]Options
}

func NewScanner(pkg *pkgx.Pkg) *Scanner { return &Scanner{pkg: pkg} }

func (s *Scanner) Options(tn *types.TypeName) (Options, bool) {
	if tn == nil {
		return nil, false
	}
	if options, ok := s.res[tn]; ok {
		return options, ok
	}

	pkg := s.pkg.PkgByPath(tn.Pkg().Path())
	if pkg == nil {
		return nil, false
	}

	for ident, def := range pkg.TypesInfo.Defs {
		c, ok := def.(*types.Const)
		if !ok {
			continue
		}
		if c.Type() != tn.Type() {
			continue
		}
		name := c.Name()
		if strings.HasPrefix(name, "_") {
			continue
		}
		spec := ident.Obj.Decl.(*ast.ValueSpec)
		label := strings.TrimSpace(spec.Comment.Text())
		val := c.Val()

		switch val.Kind() {
		case constant.String:
			v, _ := strconv.Unquote(val.String())
			s.Append(tn, NewStringOption(v, label))
		case constant.Float:
			v, _ := strconv.ParseFloat(val.String(), 64)
			s.Append(tn, NewFloatOption(v, label))
		case constant.Int:
			// TYPE_NAME_UNKNOWN
			// TYPE_NAME__XXX
			if strings.HasPrefix(name, stringsx.UpperSnakeCase(tn.Name())) {
				parts := strings.SplitN(name, "__", 2)
				if len(parts) == 2 {
					v, _ := strconv.ParseInt(val.String(), 10, 64)
					s.Append(tn, NewOption(v, parts[1], label))
				}
			} else {
				v, _ := strconv.ParseInt(val.String(), 10, 64)
				s.Append(tn, NewIntOption(v, label))
			}
		default:
			return nil, false
		}
	}
	return s.res[tn], len(s.res[tn]) > 0
}

func (s *Scanner) Append(tn *types.TypeName, opt *Option) {
	if s.res == nil {
		s.res = make(map[*types.TypeName]Options)
	}
	s.res[tn] = append(s.res[tn], *opt)
	sort.Sort(s.res[tn])
}

func (s *Scanner) Package() *pkgx.Pkg { return s.pkg }
