package client

import (
	"context"
	"sort"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
)

func NewClientGen(serviceName string, tg *TypeGen, f *g.File) *ClientGen {
	return &ClientGen{
		ServiceName: serviceName,
		f:           f,
		tg:          tg,
	}
}

type ClientGen struct {
	ServiceName string
	f           *g.File
	tg          *TypeGen
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
		SnippetClientInterface(cg.f, cg.mths...),
		SnippetNewClient(cg.f),
		SnippetClientDefine(cg.f),
		SnippetClientContextMethod(cg.f),
		SnippetClientWithContextMethod(cg.f),
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
		// rt := cg.tg.TypeInfoBySchema(ctx, mt.Schema)
		// if rt != nil {
		// 	rets = append(rets, rt.TypeInfoVar.Snippet(cg.f)...)
		// }
	}

	rets = append(
		rets,
		g.Var(g.Type(cg.f.Use(PkgKit, "Metadata"))),
		g.Var(g.Error),
	)

	m := g.Func(args...).Return(rets...).Named(op.OperationId).
		MethodOf(MethodStarReceiver("Client", "c"))

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
