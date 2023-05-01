package client

import g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"

func MethodReceiver(typ string, rcv string) *g.SnippetField {
	if rcv == "" {
		return g.Var(g.Type(typ))
	}
	return g.Var(g.Type(typ), rcv)
}

func MethodStarReceiver(typ string, rcv string) *g.SnippetField {
	if rcv == "" {
		return g.Var(g.Star(g.Type(typ)))
	}
	return g.Var(g.Star(g.Type(typ)), rcv)
}

func SnippetReturnListOfInvokeMethod(f *g.File, rt g.SnippetType) []*g.SnippetField {
	lst := make([]*g.SnippetField, 0, 3)
	if rt != nil {
		lst = append(lst, g.Var(g.Star(rt)))
	}
	lst = append(lst,
		g.Var(g.Type(f.Use(PkgKit, "Metadata"))),
		g.Var(g.Error),
	)
	return lst
}

func SnippetClientInterface(f *g.File, fns ...g.IfCanBeIfMethod) g.Snippet {
	varContext := g.Var(g.Type(f.Use("context", "Context")))
	methods := []g.IfCanBeIfMethod{
		g.Func().Named("Context").Return(varContext),
		g.Func(g.Var(g.Type(f.Use("context", "Context")))).Named("WithContext").
			Return(g.Var(g.Type("Interface"))),
	}
	for _, fn := range fns {
		methods = append(methods, fn.(*g.FuncType).WithoutBlock().WithoutReceiver())
	}

	return g.DeclType(g.Var(g.Interface(methods...), "Interface"))
}

func SnippetNewClient(f *g.File) g.Snippet {
	return g.Func(g.Var(g.Type(f.Use(PkgKit, "Client")), "c")).
		Named("NewClient").
		Return(g.Var(g.Star(g.Type("Client")))).
		Do(
			g.Return(g.Addr(g.Paren(g.Compose(
				g.Type("Client"),
				g.KeyValue(g.Ident("Client"), g.Ident("c")),
			)))),
		)
}

func SnippetClientDefine(f *g.File) g.Snippet {
	return g.DeclType(
		g.Var(g.Struct(
			g.Var(g.Type(f.Use(PkgKit, "Client")), "Client"),
			g.Var(g.Type(f.Use("context", "Context")), "ctx"),
		), "Client"),
	)
}

func SnippetClientContextMethod(f *g.File) g.Snippet {
	return g.Func().Named("Context").
		MethodOf(MethodStarReceiver("Client", "c")).
		Return(g.Var(g.Type(f.Use("context", "Context")))).
		Do(
			g.If(g.Literal(`c.ctx != nil`)).
				Do(g.Return(g.Literal("c.ctx"))),
			g.Return(g.Call(f.Use("context", "Background"))),
		)
}

func SnippetClientWithContextMethod(f *g.File) g.Snippet {
	return g.Func(g.Var(g.Type(f.Use("context", "Context")), "ctx")).
		MethodOf(MethodStarReceiver("Client", "c")).
		Named("WithContext").
		Return(g.Var(g.Type("Interface"))).
		Do(
			g.Define(g.Ident("cc")).By(g.Call("new", g.Type("Client"))),
			g.Literal("cc.Client, cc.ctx = c.Client, ctx"),
			g.Return(g.Ident("cc")),
		)
}

func SnippetOperationDefine(op string, fields ...*g.SnippetField) g.Snippet {
	return g.DeclType(g.Var(g.Struct(fields...), op))
}

func SnippetOperationPathMethod(f *g.File, op string, path string) g.Snippet {
	return g.Func().
		Named("Path").Return(g.Var(g.String)).
		MethodOf(MethodStarReceiver(op, "o")).
		Do(g.Return(f.Value(path)))
}

func SnippetOperationMethodMethod(f *g.File, op string, mtd string) g.Snippet {
	return g.Func().
		Named("Method").Return(g.Var(g.String)).
		MethodOf(MethodStarReceiver(op, "o")).
		Do(g.Return(f.Value(mtd)))
}

func SnippetOperationDoMethod(f *g.File, service, op string, comments ...string) []g.Snippet {
	return []g.Snippet{
		func() g.Snippet {
			snippet := g.Comments(comments...)
			if snippet != nil {
				return snippet
			}
			return nil
		}(),
		g.Func(
			g.Var(g.Type(f.Use("context", "Context")), "ctx"),
			g.Var(g.Type(f.Use(PkgKit, "Client")), "cli"),
			g.Var(g.Ellipsis(g.Type(f.Use(PkgKit, "Metadata"))), "metas"),
		).
			Named("Do").
			Return(g.Var(g.Type(f.Use(PkgKit, "Result")))).
			MethodOf(g.Var(g.Star(g.Type(op)), "o")).
			Do(
				g.Assign(g.Ident("ctx")).
					By(
						g.Call(
							f.Use(PkgMetax, "ContextWith"),
							g.Ident("ctx"),
							g.Valuer("operationID"),
							g.Valuer(service+"."+op),
						),
					),
				g.Return(
					// TODO impl codegen.CallSel and codegen.EllipsisVar
					g.Call("cli.Do", g.Ident("ctx"), g.Ident("o"), g.Literal("metas...")),
				),
			),
	}
}

func SnippetOperationInvokeContextMethod(f *g.File, op string, rt g.SnippetType) g.Snippet {
	return g.Func(
		g.Var(g.Type(f.Use("context", "Context")), "ctx"),
		g.Var(g.Type(f.Use(PkgKit, "Client")), "cli"),
		g.Var(g.Ellipsis(g.Type(f.Use(PkgKit, "Metadata"))), "metas"),
	).
		Return(SnippetReturnListOfInvokeMethod(f, rt)...).
		Named("InvokeContext").
		MethodOf(MethodStarReceiver(op, "o")).
		Do(
			func() g.Snippet {
				if rt != nil {
					return g.Define(g.Ident("rsp")).By(g.Call("new", rt))
				}
				return nil
			}(),
			g.Define(g.Ident("meta"), g.Ident("err")).By(
				g.Ref(
					g.Call("cli.Do", g.Ident("ctx"), g.Ident("o"), g.Literal("metas...")),
					func() g.Snippet {
						if rt != nil {
							return g.Call("Into", g.Ident("rsp"))
						} else {
							return g.Call("Into", g.Nil)
						}
					}(),
				),
			),
			func() g.Snippet {
				if rt == nil {
					return g.Return(g.Ident("meta"), g.Ident("err"))
				} else {
					return g.Return(g.Ident("rsp"), g.Ident("meta"), g.Ident("err"))
				}
			}(),
		)
}

func SnippetOperationInvokeMethod(f *g.File, op string, rt g.SnippetType) g.Snippet {
	return g.Func(
		g.Var(g.Type(f.Use(PkgKit, "Client")), "cli"),
		g.Var(g.Ellipsis(g.Type(f.Use(PkgKit, "Metadata"))), "metas"),
	).
		Return(SnippetReturnListOfInvokeMethod(f, rt)...).
		Named("Invoke").
		MethodOf(MethodStarReceiver(op, "o")).
		Do(
			// return req.InvokeContext(context.Background(), c, metas...)
			g.Return(
				g.Call(
					"o.InvokeContext",
					g.Call(f.Use("context", "Background")),
					g.Ident("cli"), g.Ident("metas..."),
				),
			),
		)
}
