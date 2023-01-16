package httpswaggergen_test

import (
	"reflect"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	. "github.com/onsi/gomega"
)

var (
	pkgImportPathHttpTransport = pkgx.Import(reflect.TypeOf(httptransport.HttpRouteMeta{}).PkgPath())
	pkgImportPathHttpx         = pkgx.Import(reflect.TypeOf(httpx.Response{}).PkgPath())
	pkgImportPathKit           = pkgx.Import(reflect.TypeOf(kit.Router{}).PkgPath())
)

func Test(t *testing.T) {
	NewWithT(t).Expect(pkgImportPathKit).To(Equal("github.com/machinefi/w3bstream/pkg/depends/kit/kit"))
	NewWithT(t).Expect(pkgImportPathHttpTransport).To(Equal("github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"))
	NewWithT(t).Expect(pkgImportPathHttpx).To(Equal("github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"))
}
