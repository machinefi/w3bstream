package openapi_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func TestOpenAPIGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "../testdata/server/cmd/app")

	ctx := log.WithLogger(context.Background(), log.Std())

	pkg, err := pkgx.LoadFrom(dir)
	NewWithT(t).Expect(err).To(BeNil())

	g := openapi.NewGenerator(pkg)

	g.Scan(ctx)
	g.Output(dir)
}
