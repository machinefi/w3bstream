package strfmtgen_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/strfmtgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func TestGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "./testdata/")

	pkg, err := pkgx.LoadFrom(dir)
	NewWithT(t).Expect(err).To(BeNil())

	file := filepath.Join(dir, "strfmt.go")

	strfmtgen.NewGenerator(pkg, file).Output(file)
}
