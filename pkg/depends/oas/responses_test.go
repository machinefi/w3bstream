package oas_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestResponse(t *testing.T) {
	g := NewCaseGroup("Response")

	g.It("empty", `{"description":""}`, Response{})
	g.It(
		"with header and content and link",
		`{"description":"desc","headers":{"x-next":{"schema":{"type":"string"}}},"content":{"application/json":{"schema":{"type":"string"}}},"links":{"GetUserByUserId":{"operationId":"getByUserId","parameters":{"userId":"$response.body#/id"}}}}`,
		func() *Response {
			resp := NewResponse("desc")
			resp.AddHeader("x-next", NewHeaderWithSchema(String()))
			resp.AddContent("application/json", NewMediaTypeWithSchema(String()))

			link := NewLink("getByUserId")
			link.AddParameter("userId", "$response.body#/id")

			resp.AddLink("GetUserByUserId", link)
			return resp
		}(),
	)

	g.Run(t)
}
