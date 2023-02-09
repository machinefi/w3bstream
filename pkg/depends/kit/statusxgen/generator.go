package statusxgen

import (
	"go/types"
	"log"
	"path"
	"path/filepath"
	"runtime"

	gen "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func New(pkg *pkgx.Pkg) *Generator {
	return &Generator{
		pkg:     pkg,
		scanner: NewScanner(pkg),
		errs:    map[string]*StatusError{},
	}
}

type Generator struct {
	pkg     *pkgx.Pkg
	scanner *Scanner
	errs    map[string]*StatusError
}

func (g *Generator) Scan(names ...string) {
	for _, name := range names {
		typeName := g.pkg.TypeName(name)
		g.errs[name] = &StatusError{
			TypeName: typeName,
			Errors:   g.scanner.StatusError(typeName),
		}
	}
}

func (g *Generator) Output(cwd string) {
	for _, e := range g.errs {
		dir, _ := filepath.Rel(
			cwd,
			must.String(pkgx.PkgPathByPath(e.TypeName.Pkg().Path())),
		)
		filename := gen.GenerateFileSuffix(path.Join(dir, stringsx.LowerSnakeCase(e.Name())+".go"))

		f := gen.NewFile(e.TypeName.Pkg().Name(), filename)

		if err := e.WriteToFile(f); err != nil {
			log.Printf("%s generate failed: %v", filename, err)
		}
	}
}

type StatusError struct {
	TypeName *types.TypeName
	Errors   []*statusx.StatusErr
}

func (s *StatusError) Name() string { return s.TypeName.Name() }

func (s *StatusError) SnippetTypeAssert(f *gen.File) gen.Snippet {
	return gen.Exprer(
		"var _ ? = (*?)(nil)",
		gen.Type(f.Use(PkgPath, "Error")),
		gen.Type(s.Name()),
	)
}

func (s *StatusError) SnippetStatusErr(f *gen.File) gen.Snippet {
	t := gen.Type(f.Use(PkgPath, "StatusErr"))

	return gen.Func().Named("StatusErr").MethodOf(gen.Var(gen.Type(s.Name()), "v")).
		Return(gen.Var(gen.Star(t))).
		Do(
			f.Expr(`return &?{
Key: v.Key(),
Code: v.Code(),
Msg: v.Msg(),
CanBeTalk: v.CanBeTalk(),
}`,
				t),
		)
}

func (s *StatusError) SnippetUnwrap(f *gen.File) gen.Snippet {
	return gen.Func().Named("Unwrap").MethodOf(gen.Var(gen.Type(s.Name()), "v")).
		Return(gen.Var(gen.Error)).
		Do(f.Expr(`return v.StatusErr()`))
}

func (s *StatusError) SnippetError(f *gen.File) gen.Snippet {
	return gen.Func().Named("Error").MethodOf(gen.Var(gen.Type(s.Name()), "v")).
		Return(gen.Var(gen.String)).
		Do(f.Expr(`return v.StatusErr().Error()`))
}

func (s *StatusError) SnippetStatusCode(f *gen.File) gen.Snippet {
	return gen.Func().Named("StatusCode").MethodOf(gen.Var(gen.Type(s.Name()), "v")).
		Return(gen.Var(gen.Int)).
		Do(
			f.Expr(
				`return ?(int(v))`,
				gen.Ident(f.Use(PkgPath, "StatusCodeFromCode")),
			),
		)
}

func (s *StatusError) SnippetCode(f *gen.File) gen.Snippet {
	return gen.Func().Named("Code").MethodOf(gen.Var(gen.Type(s.Name()), "v")).
		Return(gen.Var(gen.Int)).
		Do(
			f.Expr(`if with, ok := (interface{})(v).(?); ok {
return with.ServiceCode() + int(v)
}
return int(v)`,
				gen.Ident(f.Use(PkgPath, "ServiceCode"))),
		)
}

func (s *StatusError) SnippetKey(f *gen.File) gen.Snippet {
	clauses := make([]*gen.SnippetCaseClause, 0)

	for _, e := range s.Errors {
		clauses = append(clauses,
			gen.CaseClause(gen.Ident(e.Key)).Do(gen.Return(f.Value(e.Key))),
		)
	}

	return gen.Func().Named("Key").MethodOf(gen.Var(gen.Type(s.Name()), "v")).
		Return(gen.Var(gen.String)).
		Do(
			gen.Switch(gen.Ident("v")).
				When(
					clauses...,
				),
			gen.Return(f.Value("UNKNOWN")),
		)
}

func (s *StatusError) SnippetMsg(f *gen.File) gen.Snippet {
	clauses := make([]*gen.SnippetCaseClause, 0)

	for _, e := range s.Errors {
		clauses = append(clauses,
			gen.CaseClause(gen.Ident(e.Key)).Do(gen.Return(f.Value(e.Msg))))
	}

	return gen.Func().Named("Msg").MethodOf(gen.Var(gen.Type(s.Name()), "v")).
		Return(gen.Var(gen.String)).
		Do(
			gen.Switch(gen.Ident("v")).
				When(
					clauses...,
				),
			gen.Return(f.Value("-")),
		)
}

func (s *StatusError) SnippetCanBeTalk(f *gen.File) gen.Snippet {
	clauses := make([]*gen.SnippetCaseClause, 0)

	for _, e := range s.Errors {
		clauses = append(clauses,
			gen.CaseClause(gen.Ident(e.Key)).Do(gen.Return(f.Value(e.CanBeTalk))),
		)
	}

	return gen.Func().Named("CanBeTalk").MethodOf(gen.Var(gen.Type(s.Name()), "v")).
		Return(gen.Var(gen.Bool)).Do(
		gen.Switch(gen.Ident("v")).When(
			clauses...,
		),
		gen.Return(f.Value(false)),
	)
}

func (s *StatusError) WriteToFile(f *gen.File) error {
	f.WriteSnippet(
		s.SnippetTypeAssert(f),
		s.SnippetStatusErr(f),
		s.SnippetUnwrap(f),
		s.SnippetError(f),
		s.SnippetStatusCode(f),
		s.SnippetCode(f),
		s.SnippetKey(f),
		s.SnippetMsg(f),
		s.SnippetCanBeTalk(f),
	)
	_, err := f.Write()
	return err
}

var PkgPath string

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(current), "../statusx")
	PkgPath = must.String(pkgx.PkgIdByPath(dir))
}
