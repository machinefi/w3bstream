package httpswaggergen

import (
	"context"
	"go/ast"
	"go/types"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
	"github.com/pkg/errors"
)

func NewDefinitionScanner(pkg *pkgx.Pkg) *DefinitionScanner {
	return &DefinitionScanner{
		enumScanner:       enumgen.NewScanner(pkg),
		pkg:               pkg,
		ioWriterInterface: pkgx.New(pkg.PkgByPath("io")).TypeName("Writer").Type().Underlying().(*types.Interface),
	}
}

type DefinitionScanner struct {
	pkg               *pkgx.Pkg
	enumScanner       *enumgen.Scanner
	definitions       map[*types.TypeName]*oas.Schema
	schemas           map[string]*oas.Schema
	ioWriterInterface *types.Interface
}

func addExtension(s *oas.Schema, key string, v interface{}) {
	if s == nil {
		return
	}
	if len(s.AllOf) > 0 {
		s.AllOf[len(s.AllOf)-1].AddExtension(key, v)
	} else {
		s.AddExtension(key, v)
	}
}

func setMetaFromDoc(s *oas.Schema, doc string) {
	if s == nil {
		return
	}

	lines := strings.Split(doc, "\n")

	for i := range lines {
		if strings.Contains(lines[i], "@deprecated") {
			s.Deprecated = true
		}
	}

	description := dropMarkedLines(lines)

	if len(s.AllOf) > 0 {
		s.AllOf[len(s.AllOf)-1].Description = description
	} else {
		s.Description = description
	}
}

func fullTypeName(typeName *types.TypeName) string {
	pkg := typeName.Pkg()
	if pkg != nil {
		return pkg.Path() + "." + typeName.Name()
	}
	return typeName.Name()
}

func (scanner *DefinitionScanner) BindSchemas(openapi *oas.OpenAPI) {
	openapi.Components.Schemas = scanner.schemas
}

func (scanner *DefinitionScanner) Def(ctx context.Context, typeName *types.TypeName) *oas.Schema {
	if s, ok := scanner.definitions[typeName]; ok {
		return s
	}

	log.FromContext(ctx).Debug("scanning Type `%s.%s`", typeName.Pkg().Path(), typeName.Name())

	if typeName.IsAlias() {
		typeName = typeName.Type().(*types.Named).Obj()
	}

	doc := scanner.pkg.CommentsOf(scanner.pkg.IdentOf(typeName.Type().(*types.Named).Obj()))

	// register empty before scan
	// to avoid cycle
	scanner.setDef(typeName, &oas.Schema{})

	if doc, fmtName := parseStrfmt(doc); fmtName != "" {
		s := oas.NewSchema(oas.TypeString, fmtName)
		setMetaFromDoc(s, doc)
		return scanner.setDef(typeName, s)
	}

	if doc, typ := parseType(doc); typ != "" {
		s := oas.NewSchema(oas.Type(typ), "")
		setMetaFromDoc(s, doc)
		return scanner.setDef(typeName, s)
	}

	if typesx.FromGoType(types.NewPointer(typeName.Type())).Implements(typesx.FromGoType(scanner.ioWriterInterface)) {
		return scanner.setDef(typeName, oas.Binary())
	}

	if typeName.Pkg() != nil {
		if typeName.Pkg().Path() == "time" && typeName.Name() == "Time" {
			return scanner.setDef(typeName, oas.DateTime())
		}
	}

	if enumOptions, ok := scanner.enumScanner.Options(typeName); ok {
		s := oas.String()

		optionsLabels := make([]string, 0)
		enumVersionGot := false

		for _, o := range enumOptions {
			v := o.Value()

			if v == nil {
				continue
			}

			if !enumVersionGot {
				enumVersionGot = true

				switch v.(type) {
				case string:
					s = oas.String()
				case int64:
					s = oas.Integer()
				case float64:
					s = oas.Float()
				}
			}

			s.Enum = append(s.Enum, v)
			optionsLabels = append(optionsLabels, o.Label)
		}

		s.AddExtension(XEnumLabels, optionsLabels)

		return scanner.setDef(typeName, s)
	}

	s := oas.NewSchema(oas.TypeString, "")

	hasDefinedByInterface := false

	if method, ok := typesx.FromGoType(typeName.Type()).MethodByName("OpenAPISchemaType"); ok {
		results, n := scanner.pkg.FuncResultsOf(method.(*typesx.GoMethod).Func)
		if n == 1 {
			for _, v := range results[0] {
				if compositeLit, ok := v.Expr.(*ast.CompositeLit); ok {
					if _, ok := compositeLit.Type.(*ast.ArrayType); ok && len(compositeLit.Elts) > 0 {
						if b, ok := compositeLit.Elts[0].(*ast.BasicLit); ok {
							s.Type = oas.Type(strings.Trim(b.Value, `"`))
							hasDefinedByInterface = true
						}
					}
				}
			}
		}
	}

	if method, ok := typesx.FromGoType(typeName.Type()).MethodByName("OpenAPISchemaFormat"); ok {
		results, n := scanner.pkg.FuncResultsOf(method.(*typesx.GoMethod).Func)
		if n == 1 {
			for _, v := range results[0] {
				s.Format = strings.Trim(v.Value.String(), `"`)
				hasDefinedByInterface = true
			}
		}
	}

	if !hasDefinedByInterface {
		s = scanner.GetSchemaByType(ctx, typeName.Type().Underlying())
	}

	setMetaFromDoc(s, doc)

	return scanner.setDef(typeName, s)
}

func (scanner *DefinitionScanner) isInternal(typeName *types.TypeName) bool {
	return strings.HasPrefix(typeName.Pkg().Path(), scanner.pkg.PkgPath)
}

func (scanner *DefinitionScanner) typeUniqueName(typeName *types.TypeName, isExist func(name string) bool) (string, bool) {
	typePkgPath := typeName.Pkg().Path()
	name := typeName.Name()

	if scanner.isInternal(typeName) {
		pathParts := strings.Split(typePkgPath, "/")
		count := 1
		for isExist(name) {
			name = stringsx.UpperCamelCase(pathParts[len(pathParts)-count]) + name
			count++
		}
		return name, true
	}

	return stringsx.UpperCamelCase(typePkgPath) + name, false
}

func (scanner *DefinitionScanner) reformatSchemas() {
	typeNameList := make([]*types.TypeName, 0)

	for typeName := range scanner.definitions {
		v := typeName
		typeNameList = append(typeNameList, v)
	}

	sort.Slice(typeNameList, func(i, j int) bool {
		return scanner.isInternal(typeNameList[i]) && fullTypeName(typeNameList[i]) < fullTypeName(typeNameList[j])
	})

	schemas := map[string]*oas.Schema{}

	for _, typeName := range typeNameList {
		name, isInternal := scanner.typeUniqueName(typeName, func(name string) bool {
			_, exists := schemas[name]
			return exists
		})

		s := scanner.definitions[typeName]
		addExtension(s, XID, name)
		if !isInternal {
			addExtension(s, XGoVendorType, fullTypeName(typeName))
		}
		schemas[name] = s
	}

	scanner.schemas = schemas
}

func (scanner *DefinitionScanner) setDef(typeName *types.TypeName, schema *oas.Schema) *oas.Schema {
	if scanner.definitions == nil {
		scanner.definitions = map[*types.TypeName]*oas.Schema{}
	}
	scanner.definitions[typeName] = schema
	scanner.reformatSchemas()
	return schema
}

func NewSchemaRefer(s *oas.Schema) *SchemaRefer {
	return &SchemaRefer{
		Schema: s,
	}
}

type SchemaRefer struct {
	*oas.Schema
}

func (r SchemaRefer) RefString() string {
	s := r.Schema
	if r.Schema.AllOf != nil {
		s = r.AllOf[len(r.Schema.AllOf)-1]
	}
	return oas.NewComponentRefer("schemas", s.Extensions[XID].(string)).RefString()
}

func (scanner *DefinitionScanner) GetSchemaByType(ctx context.Context, typ types.Type) *oas.Schema {
	switch t := typ.(type) {
	case *types.Named:
		if t.String() == "mime/multipart.FileHeader" {
			return oas.Binary()
		}
		return oas.RefSchemaByRefer(NewSchemaRefer(scanner.Def(ctx, t.Obj())))
	case *types.Interface:
		return &oas.Schema{}
	case *types.Basic:
		typeName, format := getSchemaTypeFromBasicType(typesx.FromGoType(t).Kind().String())
		if typeName != "" {
			return oas.NewSchema(typeName, format)
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

		s := scanner.GetSchemaByType(ctx, elem)
		markPointer(s, count)
		return s
	case *types.Map:
		keySchema := scanner.GetSchemaByType(ctx, t.Key())
		if keySchema != nil && len(keySchema.Type) > 0 && keySchema.Type != "string" {
			panic(errors.New("only support map[string]interface{}"))
		}
		return oas.KeyValueOf(keySchema, scanner.GetSchemaByType(ctx, t.Elem()))
	case *types.Slice:
		return oas.ItemsOf(scanner.GetSchemaByType(ctx, t.Elem()))
	case *types.Array:
		length := uint64(t.Len())
		s := oas.ItemsOf(scanner.GetSchemaByType(ctx, t.Elem()))
		s.MaxItems = &length
		s.MinItems = &length
		return s
	case *types.Struct:
		structSchema := oas.ObjectOf(nil)
		schemas := make([]*oas.Schema, 0)

		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)

			if !field.Exported() {
				continue
			}

			structFieldType := field.Type()

			tags := reflect.StructTag(t.Tag(i))

			tagValueForName := tags.Get("json")
			if tagValueForName == "" {
				tagValueForName = tags.Get("name")
			}

			name, flags := tagValueAndFlagsByTagString(tagValueForName)
			if name == "-" {
				continue
			}

			if name == "" && field.Anonymous() {
				if field.Type().String() == "bytes.Buffer" {
					structSchema = oas.Binary()
					break
				}
				s := scanner.GetSchemaByType(ctx, structFieldType)
				if s != nil {
					schemas = append(schemas, s)
				}
				continue
			}

			if name == "" {
				name = field.Name()
			}

			required := true
			if hasOmitempty, ok := flags["omitempty"]; ok {
				required = !hasOmitempty
			}

			structSchema.SetProperty(
				name,
				scanner.propSchemaByField(ctx, field.Name(), structFieldType, tags, name, flags, scanner.pkg.CommentsOf(scanner.pkg.IdentOf(field))),
				required,
			)
		}

		if len(schemas) > 0 {
			return oas.AllOf(append(schemas, structSchema)...)
		}

		return structSchema
	}
	return nil
}

func (scanner *DefinitionScanner) propSchemaByField(
	ctx context.Context,
	fieldName string,
	fieldType types.Type,
	tags reflect.StructTag,
	name string,
	flags map[string]bool,
	desc string,
) *oas.Schema {
	propSchema := scanner.GetSchemaByType(ctx, fieldType)

	refSchema := (*oas.Schema)(nil)

	if propSchema.Refer != nil {
		refSchema = propSchema
		propSchema = &oas.Schema{}
		propSchema.Extensions = refSchema.Extensions
	}

	defaultValue := tags.Get("default")
	validate, hasValidate := tags.Lookup("validate")

	if flags != nil && flags["string"] {
		propSchema.Type = oas.TypeString
	}

	if defaultValue != "" {
		propSchema.Default = defaultValue
	}

	if hasValidate {
		if err := BindSchemaValidationByValidateBytes(propSchema, fieldType, []byte(validate)); err != nil {
			panic(err)
		}
	}

	setMetaFromDoc(propSchema, desc)
	propSchema.AddExtension(XGoFieldName, fieldName)

	tagKeys := map[string]string{
		"name":     XTagName,
		"mime":     XTagMime,
		"json":     XTagJSON,
		"xml":      XTagXML,
		"validate": XTagValidate,
	}

	for k, extKey := range tagKeys {
		if v, ok := tags.Lookup(k); ok {
			propSchema.AddExtension(extKey, v)
		}
	}

	if refSchema != nil {
		return oas.AllOf(
			refSchema,
			propSchema,
		)
	}

	return propSchema
}

type VendorExtensible interface {
	AddExtension(key string, value interface{})
}

func markPointer(vendorExtensible VendorExtensible, count int) {
	vendorExtensible.AddExtension(XGoStarLevel, count)
}

var (
	reStrFmt = regexp.MustCompile(`open-?api:strfmt\s+(\S+)([\s\S]+)?$`)
	reType   = regexp.MustCompile(`open-?api:type\s+(\S+)([\s\S]+)?$`)
)

func parseStrfmt(doc string) (string, string) {
	matched := reStrFmt.FindAllStringSubmatch(doc, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return doc, ""
}

func parseType(doc string) (string, string) {
	matched := reType.FindAllStringSubmatch(doc, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return doc, ""
}

var basicTypeToSchemaType = map[string][2]string{
	"invalid": {"null", ""},

	"bool":    {"boolean", ""},
	"error":   {"string", "string"},
	"float32": {"number", "float"},
	"float64": {"number", "double"},

	"int":   {"integer", "int32"},
	"int8":  {"integer", "int8"},
	"int16": {"integer", "int16"},
	"int32": {"integer", "int32"},
	"int64": {"integer", "int64"},

	"rune": {"integer", "int32"},

	"uint":   {"integer", "uint32"},
	"uint8":  {"integer", "uint8"},
	"uint16": {"integer", "uint16"},
	"uint32": {"integer", "uint32"},
	"uint64": {"integer", "uint64"},

	"byte": {"integer", "uint8"},

	"string": {"string", ""},
}

func getSchemaTypeFromBasicType(basicTypeName string) (typ oas.Type, format string) {
	if schemaTypeAndFormat, ok := basicTypeToSchemaType[basicTypeName]; ok {
		return oas.Type(schemaTypeAndFormat[0]), schemaTypeAndFormat[1]
	}
	panic(errors.Errorf("unsupported type %q", basicTypeName))
}
