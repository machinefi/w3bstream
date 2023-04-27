package client

import (
	"context"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
)

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

func NewOperationGen(serviceName string, file *g.File) *OperationGen {
	return &OperationGen{
		ServiceName: serviceName,
		f:           file,
	}
}

type OperationGen struct {
	ServiceName string
	f           *g.File
}

func (og *OperationGen) Gen(ctx context.Context, spec *oas.OpenAPI) error {
	EachOperation(spec, func(method string, path string, op *oas.Operation) {
		if op.OperationId == "" {
			return
		}
		og.Write(ctx, method, path, op)
	})
	_, err := og.f.Write()
	return err
}

func (og *OperationGen) ID(id string) string {
	if og.ServiceName != "" {
		return og.ServiceName + "." + id
	}
	return id
}

func (og *OperationGen) Write(ctx context.Context, mtd string, path string, op *oas.Operation) {
	id := op.OperationId

	fields := make([]*g.SnippetField, 0)
	for i := range op.Parameters {
		fields = append(fields, og.ParamField(ctx, op.Parameters[i]))
	}

	if field := og.RequestBodyField(ctx, op.RequestBody); field != nil {
		fields = append(fields, field)
	}

	rt, errs := og.ResponseType(ctx, &op.Responses)

	og.f.WriteSnippet(SnippetOperationDefine(id, fields...))
	og.f.WriteSnippet(SnippetOperationPathMethod(og.f, id, path))
	og.f.WriteSnippet(SnippetOperationMethodMethod(og.f, id, mtd))
	og.f.WriteSnippet(SnippetOperationDoMethod(og.f, og.ServiceName, id, errs...)...)
	og.f.WriteSnippet(SnippetOperationInvokeContextMethod(og.f, id, rt))
	og.f.WriteSnippet(SnippetOperationInvokeMethod(og.f, id, rt))
}

func (og *OperationGen) ParamField(ctx context.Context, param *oas.Parameter) *g.SnippetField {
	field := NewTypeGen(og.ServiceName, og.f).
		FieldOf(
			ctx, param.Name,
			param.Schema,
			map[string]bool{param.Name: param.Required},
		)

	tag := `in:"` + string(param.In) + `"`
	if field.Tag != "" {
		tag = tag + " " + field.Tag
	}
	field.Tag = tag

	return field
}

func (og *OperationGen) RequestBodyField(ctx context.Context, body *oas.RequestBody) *g.SnippetField {
	mt := RequestBodyMediaType(body)
	if mt == nil {
		return nil
	}

	field := NewTypeGen(og.ServiceName, og.f).FieldOf(ctx, "Data", mt.Schema, map[string]bool{})

	tag := `in:"body"`
	if field.Tag != "" {
		tag = tag + " " + field.Tag
	}
	field.Tag = tag

	return field
}

func (og *OperationGen) ResponseType(ctx context.Context, rsps *oas.Responses) (g.SnippetType, []string) {
	mt, errs := MediaTypeAndStatusErrors(rsps)
	if mt == nil {
		return nil, nil
	}
	typ, _ := NewTypeGen(og.ServiceName, og.f).Type(ctx, mt.Schema)
	return typ, errs
}
