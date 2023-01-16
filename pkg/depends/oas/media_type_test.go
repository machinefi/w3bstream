package oas_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestMediaType(t *testing.T) {
	g := NewCaseGroup("MediaType")

	g.It("empty", `{}`, MediaType{})

	g.It("with schema", `{"schema":{"type":"string"}}`, MediaType{
		MediaTypeObject: MediaTypeObject{
			Schema: String(),
		},
	})

	g.It("with schema and example", `{"schema":{"type":"string"},"example":"some string","examples":{"some":{"value":"string","externalValue":"string"}}}`, func() *MediaType {
		m := NewMediaTypeWithSchema(String())
		m.Example = "some string"

		ex := NewExample()
		ex.Value = "string"
		ex.ExternalValue = "string"

		m.AddExample("some", ex)

		return m
	}())

	g.It("with encoding", `{"schema":{"type":"string"},"encoding":{"utf-8":{"contentType":"application/json","style":"simple"}}}`, func() *MediaType {
		m := NewMediaTypeWithSchema(String())

		e := NewEncoding()
		e.ContentType = "application/json"
		e.Style = ParameterStyleSimple

		m.AddEncoding("utf-8", e)

		return m
	}())

	g.Run(t)
}
