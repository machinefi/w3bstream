package openapi_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	. "github.com/onsi/gomega"
)

func TestOpenAPIGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "../testdata/server/cmd/app")

	ctx := log.WithLogger(context.Background(), log.Std())

	pkg, err := pkgx.LoadFrom(dir)
	NewWithT(t).Expect(err).To(BeNil())

	g := openapi.NewOpenAPIGenerator(pkg)

	g.Scan(ctx)
	g.Output(dir)
}
