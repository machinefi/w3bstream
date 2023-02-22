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

func (c *Config) SnippetCommVar(f *g.File, op string) (snippets []g.Snippet) {
	switch op {
	case "filter":
		snippets = append(snippets, g.DeclVar(
			g.Var(g.Bool, "res"),
			g.Var(g.Slice(g.Byte), "ret"),
			g.Var(g.Bool, "ok"),
		),
		)
	case "map":
		snippets = append(snippets, g.DeclVar(
			g.Var(g.Type(f.Use("models", "Customer")), "res"),
			g.Var(g.Slice(g.Byte), "ret"),
			g.Var(g.Bool, "ok"),
		),
		)
	case "groupBy":
		snippets = append(snippets, g.DeclVar(
			g.Var(g.String, "res"),
			g.Var(g.Slice(g.Byte), "ret"),
			g.Var(g.Bool, "ok"),
		),
		)
	}

	return
}

func (c *Config) SnippetInDeser(f *g.File) (snippets []g.Snippet) {
	snippets = append(snippets, g.Define(
		g.Ident("src"), g.Ident("err"),
	).By(
		// TODO input
		g.Call(f.Use("json", "Marshal"), g.Ident("input")),
	),
		g.If(g.Exprer("err != nil")).
			Do(
				g.CallWith(g.Ref(g.Ident("l"), g.Ident("Error")), g.Ident("err")),
			),
	)
	return
}

func (c *Config) SnippetInvokeWasm(f *g.File, op string) (snippets []g.Snippet) {
	var doCode g.Snippet
	switch op {
	case "filter":
		doCode = g.Return(g.Ident("res"))
	case "map":
		doCode = g.Return(g.Exprer("nil"),
			g.Call(f.Use(pkgErrors, "New"), g.Valuer("the value does not support")),
		)
	case "groupBy":
		doCode = g.Return(g.Exprer("error"))
	}

	snippets = append(snippets, g.Define(g.Var(nil, "code")).
		By(
			g.CallWith(
				g.Ref(g.Ident("ins"), g.Ident("HandleEvent")),
				g.Ident("ctx"), g.Exprer("start"), g.Ident("src"),
			),
		),

		g.If(g.Exprer("code < 0")).
			Do(
				doCode,
			))
	return
}

func (c *Config) SnippetGetData(f *g.File, op string) (snippets []g.Snippet) {
	var doCode g.Snippet
	switch op {
	case "filter":
		doCode = g.Return(g.Ident("res"))
	case "map":
		doCode = g.Return(g.Exprer("nil"),
			g.Call(f.Use(pkgErrors, "New"), g.Valuer("the value does not support")),
		)
	case "groupBy":
		doCode = g.Return(g.Exprer("error"))
	}

	snippets = append(snippets, g.Assign(
		g.Ident("ret"), g.Ident("ok"),
	).By(
		g.CallWith(g.Ref(g.Ident("ins"), g.Ident("GetResource")), g.Casting(g.Uint32, g.Ident("code"))),
	),

		g.If(g.Exprer("ok")).
			Do(
				g.CallWith(
					g.Ref(g.Ident("ins"), g.Ident("RmvResource")),
					g.Ident("ctx"), g.Casting(g.Uint32, g.Ident("code")),
				).AsDefer(),
			),

		g.If(g.Exprer("!ok")).Do(
			g.CallWith(g.Ref(g.Ident("l"), g.Ident("Error")), g.Ident("err")),

			doCode,
		))
	return
}

func (c *Config) SnippetOutDeser(f *g.File, op string) (snippets []g.Snippet) {
	switch op {
	case "filter":
		snippets = append(snippets, g.Switch(g.Call(f.Use("strings", "ToLower"), g.Ident("ret"))).
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
			))
	case "map":
		snippets = append(snippets, g.Assign(
			g.Ident("err"),
		).By(
			g.Call(f.Use("json", "Unmarshal"), g.Ident("ret"), g.Ident("res")),
		),
		)
	case "groupBy":
		snippets = append(snippets, g.Assign(
			g.Ident("res"),
		).By(
			g.Casting(g.String, g.Ident("ret")),
		),
		)
	}

	return
}

func (c *Config) SnippetReturn(op string) (snippets []g.Snippet) {
	switch op {
	case "filter":
		snippets = append(snippets, g.Return(g.Ident("res")))
	case "map":
		snippets = append(snippets, g.Return(g.Ident("res"), g.Ident("err")))
	case "groupBy":
		snippets = append(snippets, g.Return(g.Ident("res")))
	}

	return
}

func (c *Config) SnippetFilterFunc(f *g.File) g.Snippet {
	var snippets []g.Snippet
	snippets = append(snippets, c.SnippetCommVar(f, "filter")...)
	snippets = append(snippets, c.SnippetInDeser(f)...)
	snippets = append(snippets, c.SnippetInvokeWasm(f, "filter")...)
	snippets = append(snippets, c.SnippetGetData(f, "filter")...)
	snippets = append(snippets, c.SnippetOutDeser(f, "filter")...)
	snippets = append(snippets, c.SnippetReturn("filter")...)

	return g.Func(g.Var(g.Interface(), "input")).Named("FilterFunc").Return(g.Var(g.Bool)).Do(snippets...)
}

func (c *Config) SnippetMapFunc(f *g.File) g.Snippet {
	var snippets []g.Snippet
	snippets = append(snippets, c.SnippetCommVar(f, "map")...)
	snippets = append(snippets, c.SnippetInDeser(f)...)
	snippets = append(snippets, c.SnippetInvokeWasm(f, "map")...)
	snippets = append(snippets, c.SnippetGetData(f, "map")...)
	snippets = append(snippets, c.SnippetOutDeser(f, "map")...)
	snippets = append(snippets, c.SnippetReturn("map")...)

	return g.Func(g.Var(g.Type(f.Use("context", "Context")), "ctx"),
		g.Var(g.Interface(), "input")).Named("MapFunc").Return(g.Var(g.Interface()), g.Var(g.Error)).Do(snippets...)
}

func (c *Config) SnippetGroupByKey(f *g.File) g.Snippet {
	var snippets []g.Snippet
	snippets = append(snippets, c.SnippetCommVar(f, "groupBy")...)
	snippets = append(snippets, c.SnippetInDeser(f)...)
	snippets = append(snippets, c.SnippetInvokeWasm(f, "groupBy")...)
	snippets = append(snippets, c.SnippetGetData(f, "groupBy")...)
	snippets = append(snippets, c.SnippetOutDeser(f, "groupBy")...)
	snippets = append(snippets, c.SnippetReturn("groupBy")...)

	return g.Func(g.Var(g.Type(f.Use("rxgo", "Item")), "item")).
		Named("GroupByKey").Return(g.Var(g.String)).Do(snippets...)
}

func (c *Config) SnippetSink() {}
