package openapi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	. "github.com/onsi/gomega"
)

func TestStatusErrScanner(t *testing.T) {
	cwd, _ := os.Getwd()
	pkg, _ := pkgx.LoadFrom(filepath.Join(cwd, "../testdata/status_err"))

	scanner := openapi.NewStatusErrScanner(pkg)

	t.Run("ScanFromComments", func(t *testing.T) {
		errs := scanner.StatusErrorsInFunc(pkg.Func("call"))
		NewWithT(t).Expect(errs).To(HaveLen(2))
		for _, e := range errs {
			t.Log(e.Summary())
		}
	})

	t.Run("ScanAll", func(t *testing.T) {
		errs := scanner.StatusErrorsInFunc(pkg.Func("main"))
		NewWithT(t).Expect(errs).To(HaveLen(3))
		for _, e := range errs {
			t.Log(e.Summary())
		}
	})
}
