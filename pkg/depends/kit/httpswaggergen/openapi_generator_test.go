package httpswaggergen_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httpswaggergen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	. "github.com/onsi/gomega"
)

func TestOpenAPIGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "./testdata/server/cmd/app")

	ctx := log.WithLogger(context.Background(), log.Std())

	pkg, err := pkgx.LoadFrom(dir)
	NewWithT(t).Expect(err).To(BeNil())

	g := httpswaggergen.NewOpenAPIGenerator(pkg)

	g.Scan(ctx)
	g.Output(dir)
}
