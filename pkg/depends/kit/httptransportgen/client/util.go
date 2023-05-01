package client

import (
	"math"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

var (
	PkgKit   = pkgx.Import(reflect.TypeOf(kit.Metadata{}).PkgPath())
	PkgMetax = pkgx.Import(reflect.TypeOf(metax.Meta{}).PkgPath())
)

func MayPrefixDeprecated(desc string, deprecated bool) []string {
	comments := []string{desc}
	if deprecated {
		comments = append([]string{"@deprecated"}, comments...)
	}
	return comments
}

func EachOperation(spec *oas.OpenAPI, mapper func(method string, path string, op *oas.Operation)) {
	ops := map[string]struct {
		Method string
		Path   string
		*oas.Operation
	}{}

	ids := make([]string, 0)

	for path := range spec.Paths.Paths {
		item := spec.Paths.Paths[path]
		for method := range item.Operations.Operations {
			op := item.Operations.Operations[method]
			// skip reserved operations for specs
			if strings.HasPrefix(op.OperationId, "OpenAPI") ||
				strings.HasPrefix(op.OperationId, "OpenAPISpec") ||
				strings.HasPrefix(op.OperationId, "ER") {
				continue
			}
			ops[op.OperationId] = struct {
				Method string
				Path   string
				*oas.Operation
			}{
				Method:    strings.ToUpper(string(method)),
				Path:      ToColonPath(path),
				Operation: op,
			}
			ids = append(ids, op.OperationId)
		}
	}

	sort.Strings(ids)

	for _, id := range ids {
		op := ops[id]
		mapper(op.Method, op.Path, op.Operation)
	}
}

func RequestBodyMediaType(body *oas.RequestBody) *oas.MediaType {
	if body == nil {
		return nil
	}

	for ct := range body.Content {
		mt := body.Content[ct]
		return mt
	}
	return nil
}

func MediaTypeAndStatusErrors(rsps *oas.Responses) (*oas.MediaType, []string) {
	if rsps == nil {
		return nil, nil
	}
	rsp := (*oas.Response)(nil)
	errs := make([]string, 0)

	for code := range rsps.Responses {
		if IsStatusCodeOK(code) {
			rsp = rsps.Responses[code]
		} else {
			exts := rsps.Responses[code].Extensions
			if exts != nil {
				if errors, ok := exts[openapi.XStatusErrs]; ok {
					if es, ok := errors.([]interface{}); ok {
						for _, err := range es {
							errs = append(errs, err.(string))
						}
					}
				}
			}
		}
	}

	sort.Strings(errs)

	if rsp == nil {
		return nil, nil
	}

	for ct := range rsp.Content {
		mt := rsp.Content[ct]
		return mt, errs
	}

	return nil, errs
}

func IsStatusCodeOK(code int) bool {
	return code >= http.StatusOK && code < http.StatusMultipleChoices
}

var regxBraceToColon = regexp.MustCompile(`/\{([^/]+)\}`)

func ToColonPath(path string) string {
	return regxBraceToColon.ReplaceAllStringFunc(path, func(str string) string {
		name := regxBraceToColon.FindAllStringSubmatch(str, -1)[0][1]
		return "/:" + name
	})
}

func ImportPathAndExpose(schema *oas.Schema) (string, string) {
	if schema.Extensions[openapi.XGoVendorType] == nil {
		return "", ""
	}
	return pkgx.ImportPathAndExpose(schema.Extensions[openapi.XGoVendorType].(string))
}

func GetEnumOptions(schema *oas.Schema) enumgen.Options {
	vs, ok := schema.Extensions[openapi.XEnumLabels]
	if !ok {
		return nil
	}
	ls, ok := vs.([]interface{})
	if !ok {
		return nil
	}

	labels := make([]string, len(ls))
	for i, v := range ls {
		if l, ok := v.(string); ok {
			labels[i] = l
		}
	}

	opts := enumgen.Options{}
	for i, e := range schema.Enum {
		o := enumgen.Option{}
		switch v := e.(type) {
		case float64:
			if math.Floor(v) == v {
				vi := int64(v)
				o.Int = &vi
			} else {
				o.Float = &v
			}
		case string:
			o.Str = &v
		}

		if len(labels) > i {
			o.Label = labels[i]
		}
		opts = append(opts, o)
	}
	return opts
}

func MayComposedAllOf(schema *oas.Schema) *oas.Schema {
	// for named field
	if schema.AllOf != nil &&
		len(schema.AllOf) == 2 &&
		schema.AllOf[len(schema.AllOf)-1].Type != oas.TypeObject {
		next := &oas.Schema{
			Reference:    schema.AllOf[0].Reference,
			SchemaObject: schema.AllOf[1].SchemaObject,
		}
		for k, v := range schema.AllOf[1].SpecExtensions.Extensions {
			next.AddExtension(k, v)
		}
		for k, v := range schema.SpecExtensions.Extensions {
			next.AddExtension(k, v)
		}
		return next
	}
	return schema
}

func PathLevels(cwd string) []string {
	paths := make([]string, 0)

	d := cwd
	for {
		paths = append(paths, d)
		if !strings.Contains(d, "/") {
			break
		}
		d = filepath.Join(d, "../")
	}

	return paths
}
