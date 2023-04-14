package client

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func SnippetsEnumDefine(f *g.File, name string, options enumgen.Options) []g.Snippet {
	if len(options) == 0 {
		return nil
	}
	ss := make([]g.Snippet, 0)

	switch options[0].Value().(type) {
	case int64:
		ss = append(ss, g.DeclType(g.Var(g.Int64, name)))
	case float64:
		ss = append(ss, g.DeclType(g.Var(g.Float64, name)))
	case string:
		ss = append(ss, g.DeclType(g.Var(g.String, name)))
	}

	sort.Sort(options)

	ss = append(ss, g.Literal("const ("))
	for _, item := range options {
		v := item.Value()
		value := v

		switch n := v.(type) {
		case string:
			value = strconv.Quote(n)
		case float64:
			vf := v.(float64)
			v = strings.Replace(strconv.FormatFloat(vf, 'f', -1, 64), ".", "_", 1)
		}

		ss = append(ss, g.Literal(fmt.Sprintf(
			`%s__%v %s = %v // %s`,
			stringsx.UpperSnakeCase(name), v, name, value, item.Label,
		)))
	}
	ss = append(ss, g.Literal(")"))
	return ss
}

func NewTypeGen(serviceName string, file *g.File) *TypeGen {
	return &TypeGen{
		ServiceName: serviceName,
		f:           file,
		enums:       map[string]enumgen.Options{},
		aliases:     map[string]string{},
	}
}

type TypeGen struct {
	ServiceName string
	f           *g.File
	enums       map[string]enumgen.Options
	aliases     map[string]string // key:schema id; val: typename // TODO optimize type alias
}

func (tg *TypeGen) Gen(ctx context.Context, spec *oas.OpenAPI) error {
	ids := make([]string, 0)
	for id := range spec.Components.Schemas {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		s := spec.Components.Schemas[id]
		typ, ok := tg.Type(ctx, s)

		if ok {
			tg.f.WriteSnippet(g.DeclType(g.Var(typ, id).AsAlias()))
			tg.aliases[id] = string(typ.Bytes())
			continue
		}

		if len(s.Properties) == 0 && s.Type == oas.TypeObject {
			path, expose := ImportPathAndExpose(s)
			if path != "" {
				if _, ok := typ.(*g.StructType); ok {
					tn := stringsx.LowerSnakeCase(path) + "." + expose
					tg.f.Use(path, strings.TrimSuffix(id, expose))
					tg.f.WriteSnippet(g.DeclType(g.Var(g.Type(tn), id).AsAlias()))
					tg.aliases[id] = tn
					continue
				}
			}
		}

		tg.f.WriteSnippet(g.DeclType(g.Var(typ, id)))
	}

	enums := make([]string, 0)
	for id := range tg.enums {
		enums = append(enums, id)
	}
	sort.Strings(enums)

	for _, enum := range enums {
		options := tg.enums[enum]
		tg.f.WriteSnippet(SnippetsEnumDefine(tg.f, enum, options)...)
	}

	_, err := tg.f.Write()
	return err
}

func (tg *TypeGen) Type(ctx context.Context, sch *oas.Schema) (g.SnippetType, bool) {
	t, alias := tg.TypeIndirect(ctx, sch)
	if sch != nil && sch.Extensions[openapi.XGoStarLevel] != nil {
		level := int(sch.Extensions[openapi.XGoStarLevel].(float64))
		for level > 0 {
			t = g.Star(t)
			level--
		}
	}
	return t, alias
}

func (tg *TypeGen) TypeIndirect(ctx context.Context, schema *oas.Schema) (g.SnippetType, bool) {
	if schema == nil {
		return g.Interface(), false
	}

	if schema.Refer != nil {
		return g.Type(schema.Refer.(*oas.ComponentRefer).ID), true
	}

	if path, expose := ImportPathAndExpose(schema); path+expose != "" {
		imports := VendorImportsFromContext(ctx)
		if len(imports) > 0 {
			for _, p := range PathLevels(path) {
				if _, ok := imports[p]; ok {
					return g.Type(tg.f.Use(path, expose)), true
				}
			}
		} else {
			return g.Type(tg.f.Use(path, expose)), true
		}
	}

	if schema.Enum != nil {
		name := stringsx.UpperCamelCase(tg.ServiceName)
		if id, ok := schema.Extensions[openapi.XID].(string); ok {
			name = name + id
			tg.enums[name] = GetEnumOptions(schema)
			return g.Type(name), true
		}
	}

	if len(schema.AllOf) > 0 {
		if schema.AllOf[len(schema.AllOf)-1].Type == oas.TypeObject {
			return g.Struct(tg.FieldsFrom(ctx, schema)...), false
		}
		return tg.TypeIndirect(ctx, MayComposedAllOf(schema))
	}

	if schema.Type == oas.TypeObject {
		if schema.AdditionalProperties != nil {
			tpe, _ := tg.Type(ctx, schema.AdditionalProperties.Schema)
			keyTyp := g.SnippetType(g.String)
			if schema.PropertyNames != nil {
				keyTyp, _ = tg.Type(ctx, schema.PropertyNames)
			}
			return g.Map(keyTyp, tpe), false
		}
		return g.Struct(tg.FieldsFrom(ctx, schema)...), false
	}

	if schema.Type == oas.TypeArray {
		if schema.Items != nil {
			t, _ := tg.Type(ctx, schema.Items)
			if schema.MaxItems != nil &&
				schema.MinItems != nil &&
				*schema.MaxItems == *schema.MinItems {
				return g.Array(t, int(*schema.MinItems)), false
			}
			return g.Slice(t), false
		}
	}

	return tg.BasicType(string(schema.Type), schema.Format), false
}

func (tg *TypeGen) BasicType(typ string, format string) g.SnippetType {
	switch format {
	case "binary":
		return g.Type(tg.f.Use("mime/multipart", "FileHeader"))
	case "byte", "int", "int8", "int16", "int32", "int64", "rune", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
		return g.BuiltInType(format)
	case "float":
		return g.Float32
	case "double":
		return g.Float64
	default:
		switch typ {
		case "null":
			// type
			return nil
		case "integer":
			return g.Int
		case "number":
			return g.Float64
		case "boolean":
			return g.Bool
		default:
			return g.String
		}
	}
}

func (tg *TypeGen) FieldsFrom(ctx context.Context, schema *oas.Schema) (fields []*g.SnippetField) {
	final := &oas.Schema{}

	if schema.AllOf != nil {
		for _, s := range schema.AllOf {
			if s.Refer != nil {
				fields = append(fields, g.Var(g.Type(s.Refer.(*oas.ComponentRefer).ID)))
			} else {
				final = s
				break
			}
		}
	} else {
		final = schema
	}

	if final.Properties == nil {
		return
	}

	names := make([]string, 0)
	for fn := range final.Properties {
		names = append(names, fn)
	}
	sort.Strings(names)

	requires := map[string]bool{}
	for _, name := range final.Required {
		requires[name] = true
	}

	for _, name := range names {
		prop := MayComposedAllOf(final.Properties[name])
		fields = append(fields, tg.FieldOf(ctx, name, prop, requires))
	}
	return
}

func (tg *TypeGen) FieldOf(ctx context.Context, name string, prop *oas.Schema, requires map[string]bool) *g.SnippetField {
	required := requires[name]

	if len(prop.AllOf) == 2 && prop.AllOf[1].Type != oas.TypeObject {
		prop = &oas.Schema{
			Reference:      prop.AllOf[0].Reference,
			SchemaObject:   prop.AllOf[1].SchemaObject,
			SpecExtensions: prop.AllOf[1].SpecExtensions,
		}
	}

	fn := stringsx.UpperCamelCase(name)
	if prop.Extensions[openapi.XGoFieldName] != nil {
		fn = prop.Extensions[openapi.XGoFieldName].(string)
	}

	typ, _ := tg.Type(ctx, prop)
	desc := MayPrefixDeprecated(prop.Description, prop.Deprecated)
	field := g.Var(typ, fn).WithComments(desc...)
	tags := map[string][]string{}

	addTagFn := func(key string, valuesOrFlags ...string) {
		tags[key] = append(tags[key], valuesOrFlags...)
	}

	addKvTagFn := func(key string, value string) {
		addTagFn(key, value)
		if !required && !strings.Contains(value, "omitempty") {
			addTagFn(key, "omitempty")
		}
	}

	if prop.Extensions[openapi.XTagJSON] != nil {
		addKvTagFn("json", prop.Extensions[openapi.XTagJSON].(string))
	}
	if prop.Extensions[openapi.XTagName] != nil {
		addKvTagFn("name", prop.Extensions[openapi.XTagName].(string))
	}
	if prop.Extensions[openapi.XTagXML] != nil {
		addKvTagFn("xml", prop.Extensions[openapi.XTagXML].(string))
	}
	if prop.Extensions[openapi.XTagMime] != nil {
		addTagFn("mime", prop.Extensions[openapi.XTagMime].(string))
	}
	if prop.Extensions[openapi.XTagValidate] != nil {
		addTagFn("validate", prop.Extensions[openapi.XTagValidate].(string))
	}
	if prop.Default != nil {
		addTagFn("default", fmt.Sprintf("%v", prop.Default))
	}

	field = field.WithTags(tags)
	return field
}
