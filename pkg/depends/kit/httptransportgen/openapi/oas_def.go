package openapi

import (
	"context"
	"go/ast"
	"go/types"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func NewDefScanner(pkg *pkgx.Pkg) *DefScanner {
	return &DefScanner{
		pkg:        pkg,
		enums:      enumgen.NewScanner(pkg),
		ioWriterIf: pkgx.New(pkg.PkgByPath("io")).TypeName("Writer").Type().Underlying().(*types.Interface),
	}
}

type DefScanner struct {
	pkg        *pkgx.Pkg
	enums      *enumgen.Scanner
	defs       map[*types.TypeName]*oas.Schema
	schemas    map[string]*oas.Schema
	ioWriterIf *types.Interface
}

func (s *DefScanner) BindSchemas(openapi *oas.OpenAPI) {
	openapi.Components.Schemas = s.schemas
}

func (s *DefScanner) Def(ctx context.Context, tn *types.TypeName) *oas.Schema {
	if schema, ok := s.defs[tn]; ok {
		return schema
	}

	log.FromContext(ctx).Debug("scanning Type `%s.%s`", tn.Pkg().Path(), tn.Name())

	if tn.IsAlias() {
		tn = tn.Type().(*types.Named).Obj()
	}

	doc := s.pkg.CommentsOf(s.pkg.IdentOf(tn.Type().(*types.Named).Obj()))

	// register empty before scan to avoid cycle
	s.set(tn, &oas.Schema{})

	// oas:strfmt
	if doc, regex := ParseStrFmt(doc); regex != "" {
		schema := oas.NewSchema(oas.TypeString, regex)
		SetMetaFromDoc(schema, doc)
		return s.set(tn, schema)
	}

	// oas:type
	if doc, typName := ParseType(doc); typName != "" {
		schema := oas.NewSchema(oas.Type(typName), "")
		SetMetaFromDoc(schema, doc)
		return s.set(tn, schema)
	}

	// impled io.Writer as oas.Binary
	if typesx.FromGoType(types.NewPointer(tn.Type())).Implements(typesx.FromGoType(s.ioWriterIf)) {
		return s.set(tn, oas.Binary())
	}

	// datetime
	if tn.Pkg() != nil {
		if tn.Pkg().Path() == "time" && tn.Name() == "Time" {
			return s.set(tn, oas.DateTime())
		}
	}

	// enums; use kit.Enum
	if opts, ok := s.enums.Options(tn); ok {
		schema := oas.String()

		labels := make([]string, 0)
		enumVersionGot := false

		for _, o := range opts {
			v := o.Value()
			if v == nil {
				continue
			}
			if !enumVersionGot {
				enumVersionGot = true

				switch v.(type) {
				case string:
					schema = oas.String()
				case int64:
					schema = oas.Integer()
				case float64:
					schema = oas.Float()
				}
			}
			schema.Enum = append(schema.Enum, v)
			labels = append(labels, o.Label)
		}
		schema.AddExtension(XEnumLabels, labels)
		return s.set(tn, schema)
	}

	// user defined interface `OpenAPISchemaType() []string`, set oas type as result[0]
	// user defined interface `OpenAPISchemaFormat() string`, set oas format as result
	schema := oas.NewSchema(oas.TypeString, "")
	hasDefinedByInterface := false

	_, res, ok := AssertIfByMtdNameAndResCntInPkg(tn.Type(), s.pkg, "OpenAPISchemaType", 1)
	if ok {
		for _, v := range res[0] {
			if lit, ok := v.Expr.(*ast.CompositeLit); ok {
				if _, ok := lit.Type.(*ast.ArrayType); ok && len(lit.Elts) > 0 {
					if b, ok := lit.Elts[0].(*ast.BasicLit); ok {
						schema.Type = oas.Type(strings.Trim(b.Value, `"`))
						hasDefinedByInterface = true
					}
				}
			}
		}
	}
	_, res, ok = AssertIfByMtdNameAndResCntInPkg(tn.Type(), s.pkg, "OpenAPISchemaFormat", 1)
	if ok {
		for _, v := range res[0] {
			schema.Format = strings.Trim(v.Value.String(), `"`)
			hasDefinedByInterface = true
		}
	}

	if !hasDefinedByInterface {
		schema = s.GetSchemaByType(ctx, tn.Type().Underlying())
	}
	SetMetaFromDoc(schema, doc)

	return s.set(tn, schema)
}

func (s *DefScanner) GetSchemaByType(ctx context.Context, typ types.Type) *oas.Schema {
	switch t := typ.(type) {
	case *types.Named:
		if t.String() == "mime/multipart.FileHeader" {
			return oas.Binary()
		}
		return oas.RefSchemaByRefer(NewSchemaRefer(s.Def(ctx, t.Obj())))
	case *types.Interface:
		return &oas.Schema{}
	case *types.Basic:
		tpe, format := GetSchemaTypeFromBasicType(typesx.FromGoType(t).Kind().String())
		if tpe != "" {
			return oas.NewSchema(tpe, format)
		}
	case *types.Pointer:
		count := 1
		elem := t.Elem()

		for {
			if p, ok := elem.(*types.Pointer); ok {
				elem = p.Elem()
				count++
			} else {
				break
			}
		}

		schema := s.GetSchemaByType(ctx, elem)
		MarkPointer(schema, count)
		return schema
	case *types.Map:
		key := s.GetSchemaByType(ctx, t.Key())
		if key != nil && len(key.Type) > 0 && key.Type != "string" {
			panic(errors.New("only support map[string]interface{}"))
		}
		return oas.KeyValueOf(key, s.GetSchemaByType(ctx, t.Elem()))
	case *types.Slice:
		return oas.ItemsOf(s.GetSchemaByType(ctx, t.Elem()))
	case *types.Array:
		length := uint64(t.Len())
		schema := oas.ItemsOf(s.GetSchemaByType(ctx, t.Elem()))
		schema.MaxItems = &length
		schema.MinItems = &length
		return schema
	case *types.Struct:
		schema := oas.ObjectOf(nil)
		schemas := make([]*oas.Schema, 0)

		for i := 0; i < t.NumFields(); i++ {
			fi := t.Field(i)
			if !fi.Exported() {
				continue
			}
			fname := fi.Name()
			fdocs := s.pkg.CommentsOf(s.pkg.IdentOf(fi))
			ftype := fi.Type()
			ftags := reflect.StructTag(t.Tag(i))

			ftag := ftags.Get("json")
			if ftag == "" {
				ftag = ftags.Get("name")
			}

			tagv, tagf := reflectx.TagValueAndFlags(ftag)
			if tagv == "-" {
				continue
			}

			if tagv == "" && fi.Anonymous() {
				if fi.Type().String() == "bytes.Buffer" {
					schema = oas.Binary()
					break
				}
				field := s.GetSchemaByType(ctx, ftype)
				if s != nil {
					schemas = append(schemas, field)
				}
				continue
			}

			if tagv == "" {
				tagv = fname
			}

			required := true
			if hasOmitempty, ok := tagf["omitempty"]; ok {
				required = !hasOmitempty
			}
			schema.SetProperty(
				tagv,
				s.propSchemaByField(ctx, fname, ftype, ftags, tagf, fdocs),
				required,
			)
		}

		if len(schemas) > 0 {
			return oas.AllOf(append(schemas, schema)...)
		}

		return schema
	}
	return nil
}

func (s *DefScanner) propSchemaByField(
	ctx context.Context,
	fname string,
	ftype types.Type,
	ftags reflect.StructTag,
	tagf map[string]bool,
	fdocs string,
) *oas.Schema {
	schema := s.GetSchemaByType(ctx, ftype)
	ref := (*oas.Schema)(nil)

	if schema.Refer != nil {
		ref = schema
		schema = &oas.Schema{}
		schema.Extensions = ref.Extensions
	}

	dfltv := ftags.Get("default")
	vldtv, hasValidate := ftags.Lookup("validate")

	if tagf != nil && tagf["string"] {
		schema.Type = oas.TypeString
	}

	if dfltv != "" {
		schema.Default = dfltv
	}

	if hasValidate {
		if err := BindSchemaValidationByValidateBytes(schema, ftype, []byte(vldtv)); err != nil {
			panic(err)
		}
	}

	SetMetaFromDoc(schema, fdocs)
	schema.AddExtension(XGoFieldName, fname)

	for tag, xkey := range TagKeysMapping {
		if v, ok := ftags.Lookup(tag); ok {
			schema.AddExtension(xkey, v)
		}
	}

	if ref != nil {
		return oas.AllOf(ref, schema)
	}

	return schema
}

func (s *DefScanner) reformat() {
	tns := make([]*types.TypeName, 0)

	for tn := range s.defs {
		tns = append(tns, tn)
	}

	sort.Slice(tns, func(i, j int) bool {
		return IsInternalType(s.pkg.PkgPath, tns[i]) &&
			FullTypeName(tns[i]) < FullTypeName(tns[j])
	})

	schemas := map[string]*oas.Schema{}

	for _, tn := range tns {
		name, isInternal := UniqueTypeName(
			s.pkg.PkgPath, tn,
			func(name string) bool {
				_, exists := schemas[name]
				return exists
			},
		)

		sch := s.defs[tn]
		AddExtension(sch, XID, name)
		if !isInternal {
			AddExtension(sch, XGoVendorType, FullTypeName(tn))
		}
		schemas[name] = sch
	}

	s.schemas = schemas
}

func (s *DefScanner) set(tn *types.TypeName, schema *oas.Schema) *oas.Schema {
	if s.defs == nil {
		s.defs = map[*types.TypeName]*oas.Schema{}
	}
	s.defs[tn] = schema
	s.reformat()
	return schema
}

type VendorExtensible interface {
	AddExtension(key string, value interface{})
}
