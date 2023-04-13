package openapi_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
)

func TestPkgPathImports(t *testing.T) {
	t.Log(PkgPathKit)
	t.Log(PkgPathStatusErr)
	t.Log(PkgPathHttpx)
	t.Log(PkgPathHttpTspt)
}
