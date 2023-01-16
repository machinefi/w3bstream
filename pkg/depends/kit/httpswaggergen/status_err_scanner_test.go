package httpswaggergen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httpswaggergen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/onsi/gomega"
)

func TestStatusErrScanner(t *testing.T) {
	cwd, _ := os.Getwd()
	pkg, _ := pkgx.LoadFrom(filepath.Join(cwd, "./testdata/status_err_scanner"))

	scanner := httpswaggergen.NewStatusErrScanner(pkg)

	t.Run("should scan from comments", func(t *testing.T) {
		statusErrs := scanner.StatusErrorsInFunc(pkg.Func("call"))
		gomega.NewWithT(t).Expect(statusErrs).To(gomega.HaveLen(2))
	})

	t.Run("should scan all", func(t *testing.T) {
		statusErrs := scanner.StatusErrorsInFunc(pkg.Func("main"))
		gomega.NewWithT(t).Expect(statusErrs).To(gomega.HaveLen(3))
	})
}
