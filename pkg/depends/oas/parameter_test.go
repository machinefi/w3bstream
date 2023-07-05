package oas_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestParameter(t *testing.T) {
	g := NewCaseGroup("Parameter")

	g.It("empty", `{"name":"","in":""}`, Parameter{})
	g.It(
		"query parameter",
		`{"name":"key","in":"query","required":true,"schema":{"type":"string"}}`,
		QueryParameter("key", String(), true),
	)

	g.It(
		"header parameter",
		`{"name":"key","in":"header","required":true,"schema":{"type":"string"}}`,
		HeaderParameter("key", String(), true),
	)

	g.It(
		"cookie parameter",
		`{"name":"key","in":"cookie","required":true,"schema":{"type":"string"}}`,
		CookieParameter("key", String(), true),
	)

	g.It(
		"path parameter",
		`{"name":"key","in":"path","required":true,"schema":{"type":"string"}}`,
		PathParameter("key", String()),
	)

	g.Run(t)
}
