package httpswaggergen

import (
	"go/types"
	"reflect"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

const (
	XID           = "x-id"
	XGoVendorType = `x-go-vendor-type`
	XGoStarLevel  = `x-go-star-level`
	XGoFieldName  = `x-go-field-name`

	XTagValidate = `x-tag-validate`
	XTagMime     = `x-tag-mime`
	XTagJSON     = `x-tag-json`
	XTagXML      = `x-tag-xml`
	XTagName     = `x-tag-name`

	XEnumLabels = `x-enum-labels`
	// Deprecated  use XEnumLabels
	XEnumOptions = `x-enum-options`
	XStatusErrs  = `x-status-errors`
)

var (
	pkgImportPathHttpTransport = pkgx.Import(reflect.TypeOf(httptransport.HttpRouteMeta{}).PkgPath())
	pkgImportPathHttpx         = pkgx.Import(reflect.TypeOf(httpx.Response{}).PkgPath())
	pkgImportPathKit           = pkgx.Import(reflect.TypeOf(kit.Router{}).PkgPath())

	pkgPathStatusx = "github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
)

func isRouterType(typ types.Type) bool {
	return strings.HasSuffix(typ.String(), pkgImportPathKit+".Router")
}

func isHttpxResponse(typ types.Type) bool {
	return strings.HasSuffix(typ.String(), pkgImportPathHttpx+".Response")
}

func isFromHttpTransport(typ types.Type) bool {
	return strings.Contains(typ.String(), pkgImportPathHttpTransport+".")
}

func tagValueAndFlagsByTagString(tagString string) (string, map[string]bool) {
	valueAndFlags := strings.Split(tagString, ",")
	v := valueAndFlags[0]
	tagFlags := map[string]bool{}
	if len(valueAndFlags) > 1 {
		for _, flag := range valueAndFlags[1:] {
			tagFlags[flag] = true
		}
	}
	return v, tagFlags
}

func filterMarkedLines(comments []string) []string {
	lines := make([]string, 0)
	for _, line := range comments {
		if !strings.HasPrefix(line, "@") {
			lines = append(lines, line)
		}
	}
	return lines
}

func dropMarkedLines(lines []string) string {
	return strings.Join(filterMarkedLines(lines), "\n")
}
