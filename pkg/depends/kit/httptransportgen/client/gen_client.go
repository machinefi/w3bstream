package client

import (
	"context"
	"sort"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
)

func SnippetClientInterface(f *g.File, name string, fns ...g.IfCanBeIfMethod) g.Snippet {
	varContext := g.Var(g.Type(f.Use("context", "Context")))
	methods := []g.IfCanBeIfMethod{
		g.Func().Named("Context").Return(varContext),
		g.Func(g.Var(g.Type(f.Use("context", "Context")))).Named("WithContext").
			Return(g.Var(g.Type(clientInterfaceName(name)))),
	}
	for _, fn := range fns {
		methods = append(methods, fn.(*g.FuncType).WithoutBlock().WithoutReceiver())
	}

	return g.DeclType(g.Var(g.Interface(methods...), clientInterfaceName(name)))
}

func SnippetNewClient(f *g.File, name string) g.Snippet {
	return g.Func(g.Var(g.Type(f.Use(PkgKit, "Client")), "c")).
		Named("New" + clientStructName(name)).
		Return(g.Var(g.Star(g.Type(clientStructName(name))))).
		Do(
			g.Return(g.Addr(g.Paren(g.Compose(
				g.Type(clientStructName(name)),
				g.KeyValue(g.Ident("Client"), g.Ident("c")),
			)))),
		)
}

func SnippetClientDefine(f *g.File, name string) g.Snippet {
	return g.DeclType(
		g.Var(g.Struct(
			g.Var(g.Type(f.Use(PkgKit, "Client")), "Client"),
			g.Var(g.Type(f.Use("context", "Context")), "ctx"),
		),
			clientStructName(name),
		),
	)
}

func SnippetClientContextMethod(f *g.File, name string) g.Snippet {
	return g.Func().Named("Context").
		MethodOf(MethodStarReceiver(clientStructName(name), "c")).
		Return(g.Var(g.Type(f.Use("context", "Context")))).
		Do(
			g.If(g.Literal(`c.ctx != nil`)).
				Do(g.Return(g.Literal("c.ctx"))),
			g.Return(g.Call(f.Use("context", "Background"))),
		)
}

func SnippetClientWithContextMethod(f *g.File, name string) g.Snippet {
	return g.Func(g.Var(g.Type(f.Use("context", "Context")), "ctx")).
		MethodOf(MethodStarReceiver(clientStructName(name), "c")).
		Named("WithContext").
		Return(g.Var(g.Type(clientInterfaceName(name)))).
		Do(
			g.Define(g.Ident("cc")).By(g.Call("new", g.Type(clientStructName(name)))),
			g.Literal("cc.Client, cc.ctx = c.Client, ctx"),
			g.Return(g.Ident("cc")),
		)
}

func NewClientGen(serviceName string, file *g.File) *ClientGen {
	return &ClientGen{
		ServiceName: serviceName,
		f:           file,
	}
}

type ClientGen struct {
	ServiceName string
	f           *g.File
	mths        []g.IfCanBeIfMethod
}

func (cg *ClientGen) Gen(ctx context.Context, spec *oas.OpenAPI) error {
	EachOperation(spec, func(method string, path string, op *oas.Operation) {
		if op.OperationId == "" {
			return
		}
		cg.mths = append(cg.mths, cg.SnippetMethodByOperation(ctx, op))
	})

	cg.Write()

	_, err := cg.f.Write()
	return err
}

func (cg *ClientGen) Write() {
	sort.Slice(cg.mths, func(i, j int) bool {
		return *cg.mths[i].(*g.FuncType).Name < *cg.mths[j].(*g.FuncType).Name
	})

	cg.f.WriteSnippet(
		SnippetClientInterface(cg.f, cg.ServiceName, cg.mths...),
		SnippetNewClient(cg.f, cg.ServiceName),
		SnippetClientDefine(cg.f, cg.ServiceName),
		SnippetClientContextMethod(cg.f, cg.ServiceName),
		SnippetClientWithContextMethod(cg.f, cg.ServiceName),
	)

	for _, m := range cg.mths {
		cg.f.WriteSnippet(m.(*g.FuncType))
	}
}

func (cg *ClientGen) SnippetMethodByOperation(ctx context.Context, op *oas.Operation) g.IfCanBeIfMethod {
	mt, _ := MediaTypeAndStatusErrors(&op.Responses)
	hasReq := len(op.Parameters) != 0 || RequestBodyMediaType(op.RequestBody) != nil

	args := make([]*g.SnippetField, 0)
	if hasReq {
		args = append(args, g.Var(g.Star(g.Type(op.OperationId)), "req"))
	}
	args = append(args, g.Var(g.Ellipsis(g.Type(cg.f.Use(PkgKit, "Metadata"))), "metas"))

	rets := make([]*g.SnippetField, 0)

	if mt != nil {
		rt, _ := NewTypeGen(cg.ServiceName, cg.f).Type(ctx, mt.Schema)
		if rt != nil {
			rets = append(rets, g.Var(g.Star(rt)))
		}
	}

	rets = append(
		rets,
		g.Var(g.Type(cg.f.Use(PkgKit, "Metadata"))),
		g.Var(g.Error),
	)

	m := g.Func(args...).Return(rets...).Named(op.OperationId).
		MethodOf(MethodStarReceiver(clientStructName(cg.ServiceName), "c"))

	if hasReq {
		return m.Do(g.Return(
			g.Exprer("req.InvokeContext(c.Context(), c.Client, metas...)")),
		)
	}

	return m.Do(g.Return(g.Exprer(
		"(&?{}).InvokeContext(c.Context(), c.Client, metas...)",
		g.Type(op.OperationId),
	)))
}
