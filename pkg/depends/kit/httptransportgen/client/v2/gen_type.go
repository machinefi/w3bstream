package client

import (
	"context"
	"fmt"
	"sort"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
	"github.com/pkg/errors"
)

func NewTypeGen(serviceName string, file *g.File) *TypeGen {
	return &TypeGen{
		ServiceName: serviceName,
		f:           file,
		types:       map[string]*TypeInfo{},
		schemas:     map[string]*oas.Schema{},
		aliases:     map[string][]*TypeInfo{},
		defines:     map[string]*TypeInfo{},
		enums:       map[string]*TypeInfo{},
	}
}

type TypeGen struct {
	ServiceName string
	f           *g.File
	schemas     map[string]*oas.Schema
	types       map[string]*TypeInfo // key: schema id
	aliases     map[string][]*TypeInfo
	defines     map[string]*TypeInfo
	enums       map[string]*TypeInfo
}

func (tg *TypeGen) Gen(ctx context.Context, spec *oas.OpenAPI) error {
	schemas := make([]*struct {
		id string
		*oas.Schema
	}, 0, len(spec.Components.Schemas))
	for id, s := range spec.Components.Schemas {
		tg.schemas[id] = s
		schemas = append(schemas, &struct {
			id string
			*oas.Schema
		}{id: id, Schema: s})
	}

	for _, s := range schemas {
		t := tg.TypeInfo(ctx, s.id)
		// defines
		if !t.IsAlias {
			if _, ok := tg.defines[s.id]; ok {
				panic(errors.Errorf("define name conflict: %s", s.id))
			}
			tg.defines[s.id] = t
			continue
		}
		// aliases & enums
		tg.aliases[t.Expose] = append(tg.aliases[t.Expose], t)
	}

	tg.WriteAliases() // skip enums
	tg.WriteEnums()
	tg.WriteDefines()

	_, err := tg.f.Write()
	return err
}

func (tg *TypeGen) WriteAliases() {
	var duplicated = true
	for duplicated {
		duplicated = false
		for alias, types := range tg.aliases {
			if len(types) <= 1 {
				continue
			}
			duplicated = true
			delete(tg.aliases, alias)
			for _, t := range types {
				t.AliasLevel++
				alias = t.Alias()
				tg.aliases[alias] = append(tg.aliases[alias], t)
			}
			break
		}
	}

	aliases := make([]*struct {
		alias string
		*TypeInfo
	}, 0, len(tg.aliases))
	for alias, types := range tg.aliases {
		if len(types) == 0 {
			continue
		}
		aliases = append(aliases, &struct {
			alias string
			*TypeInfo
		}{alias, types[0]})
	}
	sort.Slice(aliases, func(i, j int) bool {
		return aliases[i].alias < aliases[j].alias
	})

	vars := make([]g.SnippetSpec, 0)
	for _, v := range aliases {
		if v.TypeInfoEnum != nil {
			tg.enums[v.alias] = v.TypeInfo
			continue
		}
		vars = append(vars, g.Var(
			g.Type(tg.f.Use(v.Import, v.Expose)), v.alias,
		).AsAlias())
	}
	tg.f.WriteSnippet(g.DeclType(vars...))
}

func (tg *TypeGen) WriteDefines() {
	defines := make([]*struct {
		name string
		*TypeInfo
	}, 0, len(tg.defines))
	for name, t := range tg.defines {
		defines = append(defines, &struct {
			name string
			*TypeInfo
		}{name: name, TypeInfo: t})
	}
	sort.Slice(defines, func(i, j int) bool {
		return defines[i].name < defines[j].name
	})

	for _, v := range defines {
		tg.f.WriteSnippet(v.Snippet(tg.f)...)
	}
}

func (tg *TypeGen) WriteEnums() {
	enums := make([]*struct {
		alias string
		*TypeInfo
	}, 0, len(tg.enums))
	for alias, t := range tg.enums {
		enums = append(enums, &struct {
			alias string
			*TypeInfo
		}{alias: alias, TypeInfo: t})
	}
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].alias < enums[j].alias
	})

	for _, v := range enums {
		tg.f.WriteSnippet(v.Snippet(tg.f)...)
	}
}

func (tg *TypeGen) TypeInfo(ctx context.Context, id string) *TypeInfo {
	if t, ok := tg.types[id]; ok {
		return t
	}

	s := tg.schemas[id]
	if s == nil {
		panic(errors.Errorf("schema %s is not exists", id))
	}
	t := tg.IndirectTypeInfo(ctx, id, s)
	// if nil ?
	if s != nil && s.Extensions[openapi.XGoStarLevel] != nil {
		t.Star(int(s.Extensions[openapi.XGoStarLevel].(float64)))
	}
	tg.types[id] = t
	return t
}

func (tg *TypeGen) TypeInfoBySchema(ctx context.Context, id string, s *oas.Schema) *TypeInfo {
	t := tg.IndirectTypeInfo(ctx, id, s)
	// if nil ?
	if s != nil && s.Extensions[openapi.XGoStarLevel] != nil {
		t.Star(int(s.Extensions[openapi.XGoStarLevel].(float64)))
	}
	return t
}

func (tg *TypeGen) IndirectTypeInfo(ctx context.Context, id string, s *oas.Schema) *TypeInfo {
	if s == nil {
		return NewBasicTypeInfo("", "any")
	}

	if s.Refer != nil {
		return tg.TypeInfo(ctx, s.Refer.(*oas.ComponentRefer).ID).AsAlias()
	}

	path, expose := ImportPathAndExpose(s)
	t := &TypeInfo{Import: path, Expose: expose}

	// alias some type
	if path+expose != "" {
		if imports := VendorImportsFromContext(ctx); len(imports) > 0 {
			for _, p := range PathLevels(path) {
				if _, ok := imports[p]; ok {
					t = t.AsAlias()
					break
				}
			}
		} else {
			t = t.AsAlias()
		}
		if s.Enum != nil {
			return t.AsEnum(GetEnumOptions(s))
		}
		return t
	}

	// alias enum
	if s.Enum != nil {
		if _ /*id*/, ok := s.Extensions[openapi.XID].(string); ok {
			return t.AsEnum(GetEnumOptions(s)) // a single snippet
		}
	}

	// define struct
	if len(s.AllOf) > 0 {
		if last := s.AllOf[len(s.AllOf)-1]; last.Type == oas.TypeObject {
			t.Import, t.Expose = ImportPathAndExpose(last)
			if t.Import+t.Expose == "" {
				t.Expose = id
			}
			return t.AsStruct(tg.TypeInfoStructFields(ctx, id, s)...)
		}
		return tg.IndirectTypeInfo(ctx, id, MayComposedAllOf(s))
	}

	// define struct or map
	if s.Type == oas.TypeObject {
		if s.AdditionalProperties != nil {
			vt := tg.TypeInfoBySchema(ctx, id, s.AdditionalProperties.Schema)
			kt := NewBasicTypeInfo("", "string")
			if s.PropertyNames != nil {
				kt = tg.TypeInfoBySchema(ctx, id, s.PropertyNames)
			}
			return t.AsMap(kt, vt)
		}
		t.Expose = id
		return t.AsStruct(tg.TypeInfoStructFields(ctx, id, s)...)
	}

	// define slice
	if s.Type == oas.TypeArray && s.Items != nil {
		elem := tg.TypeInfoBySchema(ctx, id, s.Items)
		size := new(int)
		if s.MaxItems != nil && s.MinItems != nil && *s.MaxItems == *s.MinItems {
			size = ptrx.Ptr(int(*s.MinItems))
		}
		return t.AsSlice(elem, size)
	}

	return NewBasicTypeInfo(string(s.Type), s.Format)
}

func (tg *TypeGen) TypeInfoStructFields(ctx context.Context, id string, schema *oas.Schema) (fields []*TypeInfo) {
	s := &oas.Schema{}

	if schema.AllOf != nil {
		for _, v := range schema.AllOf {
			if v.Refer != nil {
				field := tg.TypeInfo(ctx, v.Refer.(*oas.ComponentRefer).ID)
				fields = append(fields, field)
			} else {
				s = v
				break
			}
		}
	} else {
		s = schema
	}

	if s.Properties == nil {
		return
	}

	names := make([]string, 0)
	for name := range s.Properties {
		names = append(names, name)
	}
	sort.Strings(names)

	requires := map[string]bool{}
	for _, name := range s.Required {
		requires[name] = true
	}

	for _, name := range names {
		fs := MayComposedAllOf(s.Properties[name])
		ft := tg.FieldVar(ctx, name, fs, requires[name])
		fields = append(fields, ft)
	}
	return
}

func (tg *TypeGen) FieldVar(ctx context.Context, key string, s *oas.Schema, required bool) *TypeInfo {
	if all := s.AllOf; len(all) == 2 && all[1].Type != oas.TypeObject {
		s = &oas.Schema{
			Reference:      all[0].Reference,
			SchemaObject:   all[1].SchemaObject,
			SpecExtensions: all[1].SpecExtensions,
		}
	}

	name := stringsx.UpperCamelCase(key)
	if v := s.Extensions[openapi.XGoFieldName]; v != nil {
		name = v.(string)
	}

	t := tg.TypeInfoBySchema(ctx, "", s)
	desc := MayPrefixDeprecated(s.Description, s.Deprecated)
	tags := map[string][]string{}

	if val := s.Extensions[openapi.XTagJSON]; val != nil {
		AddKvTag(tags, "json", val.(string), required)
	}
	if val := s.Extensions[openapi.XTagName]; val != nil {
		AddKvTag(tags, "name", val.(string), required)
	}
	if val := s.Extensions[openapi.XTagXML]; val != nil {
		AddKvTag(tags, "xml", val.(string), required)
	}
	if val := s.Extensions[openapi.XTagMime]; val != nil {
		AddTag(tags, "mime", val.(string))
	}
	if val := s.Extensions[openapi.XTagValidate]; val != nil {
		AddTag(tags, "validate", val.(string))
	}
	if val := s.Default; val != nil {
		AddTag(tags, "default", fmt.Sprintf("%v", val))
	}
	return t.AsVar(name, desc, tags)
}
