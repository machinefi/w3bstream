package openapi

import (
	"go/constant"
	"go/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

const (
	XID           = "x-id"
	XGoVendorType = `x-go-vendor-type`
	XGoStarLevel  = `x-go-star-level`
	XGoFieldName  = `x-go-field-name`
	XTagValidate  = `x-tag-validate`
	XTagMime      = `x-tag-mime`
	XTagJSON      = `x-tag-json`
	XTagXML       = `x-tag-xml`
	XTagName      = `x-tag-name`
	XEnumLabels   = `x-enum-labels`
	XStatusErrs   = `x-status-errors`
)

var (
	PkgPathHttpTspt  = pkgx.Import(reflect.TypeOf(httptransport.HttpRouteMeta{}).PkgPath())
	PkgPathHttpx     = pkgx.Import(reflect.TypeOf(httpx.Response{}).PkgPath())
	PkgPathKit       = pkgx.Import(reflect.TypeOf(kit.Router{}).PkgPath())
	PkgPathStatusErr = reflect.TypeOf(statusx.StatusErr{}).PkgPath()
)

var TagKeysMapping = map[string]string{
	"name":     XTagName,
	"mime":     XTagMime,
	"json":     XTagJSON,
	"xml":      XTagXML,
	"validate": XTagValidate,
}

var ParamPosOrder = map[oas.Position]string{
	"path":   "1",
	"header": "2",
	"query":  "3",
	"cookie": "4",
}

var BasicTypeAndSchemaTypeFormatMapping = map[string][2]string{
	"invalid": {"null", ""},
	"bool":    {"boolean", ""},
	"error":   {"string", "string"},
	"float32": {"number", "float"},
	"float64": {"number", "double"},
	"int":     {"integer", "int32"},
	"int8":    {"integer", "int8"},
	"int16":   {"integer", "int16"},
	"int32":   {"integer", "int32"},
	"int64":   {"integer", "int64"},
	"rune":    {"integer", "int32"},
	"uint":    {"integer", "uint32"},
	"uint8":   {"integer", "uint8"},
	"uint16":  {"integer", "uint16"},
	"uint32":  {"integer", "uint32"},
	"uint64":  {"integer", "uint64"},
	"byte":    {"integer", "uint8"},
	"string":  {"string", ""},
}

func AddExtension(s *oas.Schema, key string, v interface{}) {
	if s == nil {
		return
	}
	if len(s.AllOf) > 0 {
		s.AllOf[len(s.AllOf)-1].AddExtension(key, v)
	} else {
		s.AddExtension(key, v)
	}
}

func isRouterType(typ types.Type) bool {
	return strings.HasSuffix(typ.String(), PkgPathKit+".Router")
}

func isHttpxResponse(typ types.Type) bool {
	return strings.HasSuffix(typ.String(), PkgPathHttpx+".Response")
}

func isFromHttpTransport(typ types.Type) bool {
	return strings.Contains(typ.String(), PkgPathHttpTspt+".")
}

func isStatusError(typ types.Type, pkg *packages.Package) bool {
	t := pkgx.New(pkg).TypeName("Error").Type().Underlying().(*types.Interface)
	return types.Implements(typ, t)
}

func GetSchemaTypeFromBasicType(basic string) (typ oas.Type, format string) {
	if schemaTypeAndFormat, ok := BasicTypeAndSchemaTypeFormatMapping[basic]; ok {
		return oas.Type(schemaTypeAndFormat[0]), schemaTypeAndFormat[1]
	}
	panic(errors.Errorf("unsupported type %q", basic))
}

func FullTypeName(tn *types.TypeName) string {
	if pkg := tn.Pkg(); pkg != nil {
		return pkg.Path() + "." + tn.Name()
	}
	return tn.Name()
}

func TagFromRelativePath(pkg string, tn *types.TypeName) string {
	tag := strings.TrimPrefix(tn.Pkg().Path(), pkg)
	return strings.TrimPrefix(tag, "/")
}

func UniqueTypeName(pkg string, tn *types.TypeName, exists func(string) bool) (string, bool) {
	path := tn.Pkg().Path()
	name := tn.Name()

	if IsInternalType(pkg, tn) {
		parts := strings.Split(path, "/")
		count := 1
		for exists(name) {
			name = stringsx.UpperCamelCase(parts[len(parts)-count]) + name
			count++
		}
		return name, true
	}
	return stringsx.UpperCamelCase(path) + name, false
}

func IsInternalType(pkg string, tn *types.TypeName) bool {
	return strings.HasPrefix(tn.Pkg().Path(), pkg)
}

func AssertIfByMtdNameAndResCntInPkg(t types.Type, pkg *pkgx.Pkg, name string, cnt int) (typesx.Method, pkgx.Results, bool) {
	mtd, ok := typesx.FromGoType(t).MethodByName(name)
	if ok {
		res, n := pkg.FuncResultsOf(mtd.(*typesx.GoMethod).Func)
		if n == cnt {
			return mtd, res, true
		}
		return mtd, nil, false
	}
	return nil, nil, false
}

func MaybeTypesOfTypeName(tn *types.TypeName) []types.Type {
	return []types.Type{tn.Type(), types.NewPointer(tn.Type())}
}

func SetMetaFromDoc(s *oas.Schema, doc string) {
	if s == nil {
		return
	}

	lines := strings.Split(doc, "\n")

	for i := range lines {
		if strings.Contains(lines[i], "@deprecated") {
			s.Deprecated = true
		}
	}

	desc := DropMarkedLines(lines)

	if len(s.AllOf) > 0 {
		s.AllOf[len(s.AllOf)-1].Description = desc
	} else {
		s.Description = desc
	}
}

func DropMarkedLines(lines []string) string {
	return strings.Join(FilterMarkedLines(lines), "\n")
}

func FilterMarkedLines(comments []string) []string {
	lines := make([]string, 0)
	for _, line := range comments {
		if !strings.HasPrefix(line, "@") {
			lines = append(lines, line)
		}
	}
	return lines
}

func MarkPointer(ext VendorExtensible, count int) {
	ext.AddExtension(XGoStarLevel, count)
}

func ConstantValueOf(v constant.Value) interface{} {
	if v == nil {
		return nil
	}

	switch v.Kind() {
	case constant.Float:
		val, _ := strconv.ParseFloat(v.String(), 64)
		return val
	case constant.Bool:
		val, _ := strconv.ParseBool(v.String())
		return val
	case constant.String:
		val, _ := strconv.Unquote(v.String())
		return val
	case constant.Int:
		val, _ := strconv.ParseInt(v.String(), 10, 64)
		return val
	}

	return nil
}

func PickStatusErrorsFromDoc(doc string) []*statusx.StatusErr {
	errs := make([]*statusx.StatusErr, 0)
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		if line != "" {
			if se, err := statusx.ParseStatusErrSummary(line); err == nil {
				errs = append(errs, se)
			}
		}
	}
	return errs
}

var regxHttpRouterPath = regexp.MustCompile("/:([^/]+)")

func PatchRouterPath(openapiPath string, operation *oas.Operation) string {
	return regxHttpRouterPath.ReplaceAllStringFunc(openapiPath, func(str string) string {
		name := regxHttpRouterPath.FindAllStringSubmatch(str, -1)[0][1]

		for _, para := range operation.Parameters {
			if para.In == "path" && para.Name == name {
				return "/{" + name + "}"
			}
		}
		return "/0"
	})
}

var (
	reStrFmt = regexp.MustCompile(`open-?api:strfmt\s+(\S+)([\s\S]+)?$`)
	reType   = regexp.MustCompile(`open-?api:type\s+(\S+)([\s\S]+)?$`)
)

func ParseStrFmt(doc string) (string, string) {
	matched := reStrFmt.FindAllStringSubmatch(doc, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return doc, ""
}

func ParseType(doc string) (string, string) {
	matched := reType.FindAllStringSubmatch(doc, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return doc, ""
}
