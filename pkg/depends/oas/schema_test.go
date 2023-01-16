package oas_test

import (
	"regexp"
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

func TestSchema(t *testing.T) {
	g := NewCaseGroup("Schema")

	g.It("empty", `{}`, Schema{})

	g.It("integer", `{"type":"integer","format":"int32"}`, Integer())
	g.It("long", `{"type":"integer","format":"int64"}`, Long())
	g.It("float", `{"type":"number","format":"float"}`, Float())
	g.It("double", `{"type":"number","format":"double"}`, Double())
	g.It(
		"string",
		`{"title":"title","type":"string","description":"desc"}`,
		String().WithTitle("title").WithDesc("desc"),
	)
	g.It("byte", `{"type":"string","format":"byte"}`, Byte())
	g.It("binary", `{"type":"string","format":"binary"}`, Binary())
	g.It("date", `{"type":"string","format":"date"}`, Date())
	g.It("date-time", `{"type":"string","format":"date-time"}`, DateTime())
	g.It("password", `{"type":"string","format":"password"}`, Password())
	g.It("boolean", `{"type":"boolean"}`, Boolean())
	g.It("array", `{"type":"array","items":{"type":"string"}}`, ItemsOf(String()))
	g.It("object",
		`{"type":"object","properties":{"key1":{"type":"string"},"key2":{"type":"string"}},"required":["key1"]}`,
		ObjectOf(
			Props{
				"key1": String(),
				"key2": String(),
			},
			"key1",
		))

	g.It("object with additional",
		`{"type":"object","additionalProperties":{"type":"string"}}`,
		MapOf(String()),
	)

	g.It("object with additional and pror",
		`{"type":"object","additionalProperties":{"type":"string"},"propertyNames":{"type":"string"}}`,
		KeyValueOf(String(), String()),
	)

	g.It("anyOf", `{"anyOf":[{"type":"string"},{"type":"boolean"}]}`, AnyOf(
		String(),
		Boolean(),
	))

	g.It("oneOf", `{"oneOf":[{"type":"string"},{"type":"boolean"}]}`, OneOf(
		String(),
		Boolean(),
	))

	g.It("allOf", `{"allOf":[{"type":"string"},{"type":"boolean"}]}`, AllOf(
		String(),
		Boolean(),
	))

	g.It("not", `{"not":{"type":"string"}}`, Not(String()))

	validation := &SchemaValidation{
		MultipleOf:       ptrx.Ptr(float64(2)),
		Maximum:          ptrx.Ptr(float64(10)),
		ExclusiveMaximum: true,
		Minimum:          ptrx.Ptr(float64(1)),
		ExclusiveMinimum: true,

		MaxLength: ptrx.Ptr(uint64(10)),
		MinLength: ptrx.Ptr(uint64(0)),
		Pattern:   regexp.MustCompile("/+d/").String(),

		MaxItems:    ptrx.Ptr(uint64(10)),
		MinItems:    ptrx.Ptr(uint64(1)),
		UniqueItems: true,

		MaxProperties: ptrx.Ptr(uint64(10)),
		MinProperties: ptrx.Ptr(uint64(1)),
		Required:      []string{"key"},

		Enum: []interface{}{"1", "2", "3"},
	}

	g.It(
		"with string validation",
		`{"type":"string","maxLength":10,"minLength":0,"pattern":"/+d/","enum":["1","2","3"]}`,
		String().WithValidation(validation),
	)

	g.It(
		"with integer validation",
		`{"type":"integer","format":"int32","multipleOf":2,"maximum":10,"exclusiveMaximum":true,"minimum":1,"exclusiveMinimum":true,"enum":["1","2","3"]}`,
		Integer().WithValidation(validation),
	)

	g.It(
		"with number validation",
		`{"type":"number","format":"float","multipleOf":2,"maximum":10,"exclusiveMaximum":true,"minimum":1,"exclusiveMinimum":true,"enum":["1","2","3"]}`,
		Float().WithValidation(validation),
	)

	g.It(
		"with array validation",
		`{"type":"array","items":{"type":"string"},"maxItems":10,"minItems":1,"uniqueItems":true,"enum":["1","2","3"]}`,
		ItemsOf(String()).WithValidation(validation),
	)

	g.It("object",
		`{"type":"object","properties":{"key1":{"type":"string"},"key2":{"type":"string"}},"maxProperties":10,"minProperties":1,"required":["key"],"enum":["1","2","3"]}`,
		ObjectOf(
			Props{
				"key1": String(),
				"key2": String(),
			},
			"key1",
		).WithValidation(validation),
	)

	g.Run(t)
}
