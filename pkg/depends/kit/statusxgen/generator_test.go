package statusxgen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusxgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func TestGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	pkg, _ := pkgx.LoadFrom(filepath.Join(cwd, "./__examples__"))

	g := statusxgen.New(pkg)

	g.Scan("StatusError")
	g.Output(cwd)
}
