package oas_test

import (
	"net/http"
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestOperation(t *testing.T) {
	g := NewCaseGroup("Operation")

	g.It("empty", `{"operationId":"","responses":{}}`, Operation{})

	g.It("full", `{"tags":["pets"],"summary":"List all pets","description":"desc","operationId":"listPets","parameters":[{"name":"limit","in":"query","description":"How many items to return at one time (max 100)","required":true,"schema":{"type":"integer","format":"int32"}},{"name":"size","in":"query","schema":{"type":"integer","format":"int32"}}],"responses":{"200":{"description":"An paged array of pets","headers":{"x-next":{"schema":{"type":"string","description":"A link to the next page of responses"}}},"content":{"application/json":{"schema":{"type":"string"}}}},"default":{"description":"unexpected error","content":{"text/html":{"schema":{"type":"string"}}}}},"callbacks":{"myEvent":{"{$request.query.callbackUrl}?event={$request.query.event}":{"post":{"operationId":"callback","requestBody":{"content":{"application/json":{"schema":{"type":"string"}}}},"responses":{"200":{"description":"OK"}}}}}}}`, func() *Operation {
		op := NewOperation("listPets").WithSummary("List all pets").WithDesc("desc").WithTags("pets")

		op.AddParameter(
			QueryParameter("limit", Integer(), true).WithDesc("How many items to return at one time (max 100)"),
		)

		op.AddParameter(
			QueryParameter("size", Integer(), false),
		)

		{
			resp := NewResponse("An paged array of pets")

			s := String()
			s.Description = "A link to the next page of responses"
			resp.AddHeader("x-next", NewHeaderWithSchema(s))
			resp.AddContent("application/json", NewMediaTypeWithSchema(String()))

			op.AddResponse(http.StatusNoContent, nil)
			op.AddResponse(http.StatusOK, resp)
		}

		{
			resp := NewResponse("unexpected error")
			resp.AddContent("text/html", NewMediaTypeWithSchema(String()))

			op.SetDefaultResponse(nil)
			op.SetDefaultResponse(resp)
		}

		opCallback := NewOperation("callback")
		rb := NewRequestBody("", false)
		rb.AddContent("application/json", NewMediaTypeWithSchema(String()))
		opCallback.SetRequestBody(rb)

		opCallback.AddResponse(http.StatusOK, NewResponse("OK"))

		op.AddCallback("myEvent", NewCallback(POST, `{$request.query.callbackUrl}?event={$request.query.event}`, opCallback))

		return op
	}())

	g.Run(t)
}
