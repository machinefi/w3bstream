package client

import (
	"context"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
)

func NewOperationGen(serviceName string, tg *TypeGen, f *g.File) *OperationGen {
	return &OperationGen{
		ServiceName: serviceName,
		f:           f,
		tg:          tg,
	}
}

type OperationGen struct {
	ServiceName string
	f           *g.File
	tg          *TypeGen
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
	// TODO
	// field := og.tg.FieldVar(ctx, param.Name, param.Schema, param.Required)

	// if field.Tags == nil {
	// 	field.Tags = make(map[string][]string)
	// }

	// AddTag(field.Tags, "in", string(param.In))
	// return field.Snippet().(*g.SnippetField)
	return nil
}

func (og *OperationGen) RequestBodyField(ctx context.Context, body *oas.RequestBody) *g.SnippetField {
	// TODO
	// mt := RequestBodyMediaType(body)
	// if mt == nil {
	// 	return nil
	// }

	// field := og.tg.FieldVar(ctx, "Data", mt.Schema, false)

	// if field.Tags == nil {
	// 	field.Tags = make(map[string][]string)
	// }

	// AddTag(field.Tags, "in", "body")
	// return field.Snippet().(*g.SnippetField)
	return nil
}

func (og *OperationGen) ResponseType(ctx context.Context, rsps *oas.Responses) (g.SnippetType, []string) {
	// return nil
	// mt, errs := MediaTypeAndStatusErrors(rsps)
	// if mt == nil {
	// 	return nil, nil
	// }
	// t := NewTypeGen(og.ServiceName, og.f).TypeInfoBySchema(ctx, mt.Schema)
	// return t.Snippet().(g.SnippetType), errs
	return nil, nil
}
