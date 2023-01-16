package openapi_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func ExampleNewRouterScanner() {
	cwd, _ := os.Getwd()
	pkg, _ := pkgx.LoadFrom(filepath.Join(cwd, "../testdata/router_scanner"))

	router := pkg.Var("Router")

	scanner := openapi.NewRouterScanner(pkg)
	routes := scanner.Router(router).Routes()

	for _, r := range routes {
		fmt.Println(r.String())
	}
	// Output:
	// GET /root/:id httptransport.MetaOperator auth.Auth main.Get
	// HEAD /root/group/health httptransport.MetaOperator httptransport.MetaOperator httptransport.MetaOperator group.Health
}
