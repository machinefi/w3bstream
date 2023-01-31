package codegen

import g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"

type Config struct{}

var (
	pkgLogger   = "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	pkgSqlx     = "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	pkgInstance = "github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
	pkgErrors   = "github.com/pkg/errors"
)

func (c *Config) SnippetMain(f *g.File) g.Snippet {
	return g.Func().
		Named("main").
		Do(
			g.DeclVar(
				g.Var(g.Type(f.Use("context", "Context")), "ctx"),
				g.Var(g.Type(f.Use(pkgLogger, "Logger")), "l"),
				g.Var(g.Type(f.Use(pkgSqlx, "DBExecutor")), "d"),
				g.Var(g.Type(f.Use(pkgInstance, "Instance")), "ins"),
				g.Var(g.Error, "err"),
			),

			// TODO
		)
}

func (c *Config) SnippetInitSource() {}

func (c *Config) SnippetInitSink() {}

func (c *Config) SnippetInitChannel() {}

func (c *Config) SnippetCompute() {}

func (c *Config) SnippetFilterFunc(f *g.File) g.Snippet {
	return g.Func(g.Var(g.Any, "input")).
		Named("FilterFunc").
		Return(g.Var(g.Bool)).
		Do(
			g.DeclVar(g.Var(g.Bool, "res")),

			g.Define(
				g.Ident("src"), g.Ident("err"),
			).By(
				g.Call(f.Use("json", "Marshal"), g.Ident("input")),
			),

			g.If(g.Exprer("err != nil")).
				Do(
					g.CallWith(g.Ref(g.Ident("l"), g.Ident("Error")), g.Ident("err")),
				),

			g.Define(g.Var(nil, "code")).
				By(
					g.CallWith(
						g.Ref(g.Ident("ins"), g.Ident("HandleEvent")),
						g.Ident("ctx"), g.Exprer("start"), g.Ident("src"),
					),
				),

			g.If(g.Exprer("code < 0")).
				Do(
					g.Return(g.Ident("res")),
				),

			g.Define(
				g.Ident("ret"), g.Ident("ok"),
			).By(
				g.CallWith(g.Ref(g.Ident("ins"), g.Ident("GetResource"), g.Casting(g.Uint32, g.Ident("code")))),
			),

			g.If(g.Exprer("ok")).
				Do(
					g.CallWith(
						g.Ref(g.Ident("ins"), g.Ident("RmvResource"),
							g.Ident("ctx"), g.Casting(g.Uint32, g.Ident("code"))),
					).AsDefer(),
				),

			g.If(g.Exprer("!ok")).Do(g.Return(g.Ident("res"))),

			g.Switch(g.Call(f.Use("strings", "ToLower"), g.Ident("ret"))).
				When(
					g.CaseClause(g.Valuer("true")).Do(
						g.Assign(g.Ident("res")).By(g.True),
					),
					g.CaseClause(g.Valuer("false")).Do(
						g.Assign(g.Ident("res")).By(g.False),
					),
					g.CaseClause().Do(
						g.CallWith(
							g.Ref(g.Ident("l"), g.Ident("Warn")),
							g.Call(f.Use(pkgErrors, "New"), g.Valuer("the value does not support"))),
					),
				),

			g.Return(g.Ident("res")),
		)
}

func (c *Config) SnippetMapFunc() {}

func (c *Config) SnippetGroupByKey() {}

func (c *Config) SnippetSink() {}
