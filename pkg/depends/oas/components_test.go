package oas_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/stretchr/testify/assert"
)

func TestComponents(t *testing.T) {
	g := NewCaseGroup("Components")

	components := &Components{}

	components.AddParameter("key", QueryParameter("key", String(), false))
	components.AddParameter("nothing", nil)

	components.AddHeader("key", NewHeaderWithSchema(String()))
	components.AddHeader("nothing", nil)

	components.AddSecurityScheme("key", NewAPIKeySecurityScheme("AccessToken", PositionHeader))
	components.AddSecurityScheme("nothing", nil)

	components.AddLink("key", NewLink("link"))
	components.AddLink("nothing", nil)

	components.AddCallback("key", NewCallback(POST, "callback", NewOperation("op")))
	components.AddCallback("nothing", nil)

	components.AddResponse("key", NewResponse("ok"))
	components.AddResponse("nothing", nil)

	components.AddSchema("key", String())
	components.AddSchema("nothing", nil)

	components.AddExample("key", NewExample())
	components.AddExample("nothing", nil)

	components.AddRequestBody("key", NewRequestBody("desc", false))
	components.AddRequestBody("nothing", nil)

	g.It("all", `{"schemas":{"key":{"type":"string"}},"responses":{"key":{"description":"ok"}},"parameters":{"key":{"name":"key","in":"query","schema":{"type":"string"}}},"examples":{"key":{}},"requestBodies":{"key":{"description":"desc"}},"headers":{"key":{"schema":{"type":"string"}}},"securitySchemes":{"key":{"type":"apiKey","name":"AccessToken","in":"header"}},"links":{"key":{"operationId":"link"}},"callbacks":{"key":{"callback":{"post":{"operationId":"op","responses":{}}}}}}`, components)

	g.It("ref parameter", `{"$ref":"#/components/parameters/key"}`, components.RefParameter("key"))
	g.It("ref header", `{"$ref":"#/components/headers/key"}`, components.RefHeader("key"))
	g.It("ref link", `{"$ref":"#/components/links/key"}`, components.RefLink("key"))
	g.It("ref callback", `{"$ref":"#/components/callbacks/key"}`, components.RefCallback("key"))
	g.It("ref response", `{"$ref":"#/components/responses/key"}`, components.RefResponse("key"))
	g.It("ref schema", `{"$ref":"#/components/schemas/key"}`, components.RefSchema("key"))
	g.It("ref example", `{"$ref":"#/components/examples/key"}`, components.RefExample("key"))
	g.It("ref request body", `{"$ref":"#/components/requestBodies/key"}`, components.RefRequestBody("key"))
	g.It("request body", `{"key":[]}`, components.RequireSecurity("key"))

	assert.Nil(t, components.RefParameter("not_found"))
	assert.Nil(t, components.RefHeader("not_found"))
	assert.Nil(t, components.RefLink("not_found"))
	assert.Nil(t, components.RefCallback("not_found"))
	assert.Nil(t, components.RefResponse("not_found"))
	assert.Nil(t, components.RefSchema("not_found"))
	assert.Nil(t, components.RefExample("not_found"))
	assert.Nil(t, components.RefRequestBody("not_found"))

	g.Run(t)
}
