package openapi_test

import (
	"encoding/json"
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
)

func TestPkgPathImports(t *testing.T) {
	t.Log(PkgPathKit)
	t.Log(PkgPathStatusErr)
	t.Log(PkgPathHttpx)
	t.Log(PkgPathHttpTspt)
}

func TestMarshalNameRules(t *testing.T) {
	content, _ := json.Marshal([]string{
		"projectName", "^[a-z0-9_]{6,32}$", "project-name", "^[a-z0-9_]{6,32}$",
	})
	t.Log(string(content))
}
