package client

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func NewBasicTypeInfo(typ, format string) *TypeInfo {
	t := &TypeInfo{}

	switch format {
	case "binary":
		t.Import, t.Expose = "mime/multipart", "FileHeader"
	case "byte", "int", "int8", "int16", "int32", "int64", "rune", "uint", "uint8", "uint16",
		"uint32", "uint64", "uintptr", "float32", "float64", "string", "interface{}", "any":
		t.Expose = format
	case "float":
		t.Expose = "float32"
	case "double":
		t.Expose = "float64"
	default:
		switch typ {
		case "null":
			return nil // TODO skip this type?
		case "integer":
			t.Expose = "int"
		case "number":
			t.Expose = "float64"
		case "boolean":
			t.Expose = "bool"
		default:
			t.Expose = "string"
		}
	}

	return t
}

type TypeInfo struct {
	Expose  string
	Import  string
	IsAlias bool
	// Example for type alias 'github.com/repo/pkg/sub1/sub2/sub3.SomeType'
	// will alias to Sub3SomeType Sub2Sub3SomeType Sub1Sub2Sub3SomeType
	// until the full ID GithubComRepoPkgSub1Sub2Sub3SomeType to resolve name conflict
	AliasLevel int

	*TypeInfoStar
	*TypeInfoEnum
	*TypeInfoStruct
	*TypeInfoMap
	*TypeInfoSlice
	*TypeInfoVar
}

func (t TypeInfo) AsAlias() *TypeInfo {
	t.IsAlias = true
	return &t
}

func (t TypeInfo) SetAliasLevel(level int) *TypeInfo {
	t.AliasLevel = level
	return &t
}

func (t TypeInfo) Star(level int) *TypeInfo {
	t.TypeInfoStar = &TypeInfoStar{Level: level}
	return &t
}

func (t TypeInfo) AsEnum(opts enumgen.Options) *TypeInfo {
	t.TypeInfoEnum = &TypeInfoEnum{Options: opts}
	return &t
}

func (t TypeInfo) AsStruct(fields ...*TypeInfo) *TypeInfo {
	t.TypeInfoStruct = &TypeInfoStruct{Fields: fields}
	return &t
}

func (t TypeInfo) AsMap(k, v *TypeInfo) *TypeInfo {
	t.TypeInfoMap = &TypeInfoMap{KeyType: k, ValType: v}
	return &t
}

func (t TypeInfo) AsSlice(elem *TypeInfo, size *int) *TypeInfo {
	t.TypeInfoSlice = &TypeInfoSlice{ElemType: elem, Size: size}
	return &t
}

func (t TypeInfo) AsVar(name string, desc []string, tags map[string][]string) *TypeInfo {
	t.TypeInfoVar = &TypeInfoVar{Name: name, Desc: desc, Tags: tags}
	return &t
}

func (t *TypeInfo) Alias() string {
	if t.AliasLevel == 0 || t.Import == "" {
		return t.Expose
	}
	routes := append(LastPackages(t.Import, t.AliasLevel), t.Expose)
	return stringsx.UpperCamelCase(strings.Join(routes, "-"))
}

func (t *TypeInfo) SnippetType() g.SnippetType {
	return g.Type(t.Alias())
}

func (t *TypeInfo) Snippet(f *g.File) (ss []g.Snippet) {
	if t.IsAlias && t.TypeInfoEnum == nil {
		ss = append(ss, g.DeclType(g.Var(
			g.Type(f.Use(t.Import, t.Expose)),
			t.Alias(),
		).AsAlias()))
		return
	}
	if t.TypeInfoEnum != nil {
		ss = append(ss, t.TypeInfoEnum.Snippet(t)...)
		return
	}
	if t.TypeInfoStruct != nil {
		ss = append(ss, g.DeclType(g.Var(g.Struct(), t.Alias())))
	}
	return ss
}

type TypeInfoStar struct{ Level int }

func (t *TypeInfoStar) Snippet(ti *TypeInfo) g.SnippetType {
	typ := ti.SnippetType()
	for level := t.Level; level > 0; level-- {
		typ = g.Star(typ)
	}
	return typ
}

type TypeInfoEnum struct{ Options enumgen.Options }

func (t *TypeInfoEnum) Snippet(ti *TypeInfo) []g.Snippet {
	if len(t.Options) == 0 {
		return nil
	}
	ss := make([]g.Snippet, 0)

	name := ti.Alias()

	switch t.Options[0].Value().(type) {
	case int64:
		ss = append(ss, g.DeclType(g.Var(g.Int64, name)))
	case float64:
		ss = append(ss, g.DeclType(g.Var(g.Float64, name)))
	case string:
		ss = append(ss, g.DeclType(g.Var(g.String, name)))
	}
	sort.Sort(t.Options)

	snippet := ""
	for _, item := range t.Options {
		v := item.Value()
		value := v

		switch n := v.(type) {
		case string:
			value = strconv.Quote(n)
		case float64:
			vf := v.(float64)
			v = strings.Replace(strconv.FormatFloat(vf, 'f', -1, 64), ".", "_", 1)
		}

		_ = fmt.Sprintf("%s__%v %s = %v // %s\n", stringsx.UpperSnakeCase(name), v, name, value, item.Label)
	}

	for _, item := range t.Options {
		label := item.Label
		value := item.Value()
		suffix := value

		switch v := value.(type) {
		case string:
			value = strconv.Quote(v)
		case float64:
			value = strconv.FormatFloat(v, 'f', -1, 64)
		case int64:
			value = strconv.FormatInt(v, 10)
		}
		snippet += fmt.Sprintf("%s__%v %s = %v // %s\n",
			stringsx.UpperSnakeCase(name), suffix, name,
			value, stringsx.UpperSnakeCase(label),
		)
	}
	ss = append(ss, g.Literal(fmt.Sprintf("const(\n%s)", snippet)))
	return ss
}

type TypeInfoStruct struct{ Fields []*TypeInfo }

type TypeInfoMap struct{ KeyType, ValType *TypeInfo }

type TypeInfoSlice struct {
	ElemType *TypeInfo
	Size     *int
}

type TypeInfoVar struct {
	Name string
	Desc []string
	Tags map[string][]string
}

func (t *TypeInfoVar) Snippet(f *g.File) []*g.SnippetField { return nil }
