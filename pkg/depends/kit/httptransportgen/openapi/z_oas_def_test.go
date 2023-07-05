package openapi_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func TestDefinitionScanner(t *testing.T) {
	cwd, _ := os.Getwd()

	pkg, err := pkgx.LoadFrom(filepath.Join(cwd, "../testdata/definition_scanner"))
	NewWithT(t).Expect(err).To(BeNil())

	scanner := openapi.NewDefScanner(pkg)

	cases := [][2]string{
		{
			"Node", // language=JSON
			`{
  "type": "object",
  "properties": {
    "children": {
      "type": "array",
      "items": {
        "$ref": "#/components/schemas/Node"
      },
      "x-go-field-name": "Children",
      "x-tag-json": "children"
    },
    "type": {
      "type": "string",
      "x-go-field-name": "Type",
      "x-tag-json": "type"
    }
  },
  "required": [
    "type",
    "children"
  ],
  "x-id": "Node"
}`},
		{
			"Interface", // language=JSON
			`{
  "x-id": "Interface"
}
`}, {
			"Binary", // language=JSON
			`{
  "type": "string",
  "format": "binary",
  "x-go-star-level": 1,
  "x-id": "Binary"
}`}, {
			"String", // language=JSON
			`{
  "type": "string",
  "x-id": "String"
}`}, {
			"String", // language=JSON
			`{
  "type": "string",
  "x-id": "String"
}`}, {
			"Bool", // language=JSON
			`{
  "type": "boolean",
  "x-id": "Bool"
}`}, {

			"Float", // language=JSON
			`{
  "type": "number",
  "format": "float",
  "x-id": "Float"
}`}, {
			"Double", // language=JSON
			`{
  "type": "number",
  "format": "double",
  "x-id": "Double"
}`}, {
			"Int", // language=JSON
			`{
  "type": "integer",
  "format": "int32",
  "x-id": "Int"
}`}, {
			"Uint", // language=JSON
			`{
  "type": "integer",
  "format": "uint32",
  "x-id": "Uint"
}`}, {
			"Time", // language=JSON
			`{
  "type": "string",
  "format": "date-time",
  "description": "日期",
  "x-id": "Time"
}
`}, {
			"FakeBool", // language=JSON
			`{
  "type": "boolean",
  "x-id": "FakeBool"
}
`}, {
			"Enum", // language=JSON
			`{
  "type": "string",
  "enum": [
    "ONE",
    "TWO"
  ],
  "x-enum-labels": [
    "one",
    "two"
  ],
  "x-id": "Enum"
}`}, {
			"Map", // language=JSON
			`{
  "type": "object",
  "additionalProperties": {
    "$ref": "#/components/schemas/String"
  },
  "propertyNames": {
    "type": "string"
  },
  "x-id": "Map"
}`}, {
			"ArrayString", // language=JSON
			`{
  "type": "array",
  "items": {
    "type": "string"
  },
  "maxItems": 2,
  "minItems": 2,
  "x-id": "ArrayString"
}`}, {
			"SliceString", // language=JSON
			`{
  "type": "array",
  "items": {
    "type": "string"
  },
  "x-id": "SliceString"
}`}, {
			"SliceNamed", // language=JSON
			`{
  "type": "array",
  "items": {
    "$ref": "#/components/schemas/String"
  },
  "x-id": "SliceNamed"
}`}, {
			"TimeAlias", // language=JSON
			`{
  "type": "string",
  "format": "date-time",
  "x-go-vendor-type": "time.Time",
  "x-id": "TimeTime"
}`}, {
			"Struct", // language=JSON
			`{
  "type": "object",
  "properties": {
    "createdAt": {
      "allOf": [
        {
          "$ref": "#/components/schemas/TimeTime"
        },
        {
          "x-go-field-name": "CreatedAt",
          "x-tag-json": "createdAt,omitempty"
        }
      ]
    },
    "enum": {
      "allOf": [
        {
          "$ref": "#/components/schemas/Enum"
        },
        {
          "type": "string",
          "enum": [
            "ONE"
          ],
          "x-go-field-name": "Enum",
          "x-tag-json": "enum",
          "x-tag-validate": "@string{ONE}"
        }
      ]
    },
    "id": {
      "type": "string",
      "minLength": 0,
      "pattern": "\\d+",
      "default": "1",
      "x-go-field-name": "ID",
      "x-go-star-level": 2,
      "x-tag-json": "id,omitempty",
      "x-tag-validate": "@string/\\d+/"
    },
    "map": {
      "type": "object",
      "additionalProperties": {
        "type": "object",
        "additionalProperties": {
          "type": "object",
          "properties": {
            "id": {
              "type": "integer",
              "format": "int32",
              "maximum": 10,
              "minimum": 0,
              "x-go-field-name": "ID",
              "x-tag-json": "id",
              "x-tag-validate": "@int[0,10]"
            }
          },
          "required": [
            "id"
          ]
        },
        "propertyNames": {
          "type": "string"
        },
        "minProperties": 0
      },
      "propertyNames": {
        "type": "string"
      },
      "maxProperties": 3,
      "minProperties": 0,
      "x-go-field-name": "Map",
      "x-tag-json": "map,omitempty",
      "x-tag-validate": "@map\u003c,@map\u003c,@struct\u003e\u003e[0,3]"
    },
    "name": {
      "type": "string",
      "minLength": 2,
      "description": "name",
      "x-go-field-name": "Name",
      "x-go-star-level": 1,
      "x-tag-json": "name",
      "x-tag-validate": "@string[2,]"
    },
    "slice": {
      "type": "array",
      "items": {
        "type": "number",
        "format": "double"
      },
      "maxItems": 3,
      "minItems": 1,
      "x-go-field-name": "Slice",
      "x-tag-json": "slice",
      "x-tag-validate": "@slice\u003c@float64\u003c7,5\u003e\u003e[1,3]"
    }
  },
  "required": [
    "name",
    "enum",
    "slice"
  ],
  "x-id": "Struct"
}`}, {
			"Composed", // language=JSON
			`{
  "allOf": [
    {
      "$ref": "#/components/schemas/Part"
    },
    {
      "type": "object",
      "description": "Composed",
      "x-id": "Composed"
    }
  ]
}`}, {
			"NamedComposed", // language=JSON
			`{
  "type": "object",
  "properties": {
    "part": {
      "allOf": [
        {
          "$ref": "#/components/schemas/Part"
        },
        {
          "x-go-field-name": "Part",
          "x-tag-json": "part"
        }
      ]
    }
  },
  "required": [
    "part"
  ],
  "x-id": "NamedComposed"
}`},
	}

	for _, c := range cases {
		t.Run(c[0], func(t *testing.T) {
			s := scanner.Def(context.Background(), pkg.TypeName(c[0]))
			data, _ := json.MarshalIndent(s, "", "  ")
			NewWithT(t).Expect(string(data)).To(Equal(strings.TrimSpace(c[1])))
		})
	}

	t.Run("bind", func(t *testing.T) {
		openAPI := oas.NewOpenAPI()
		openAPI.AddOperation(oas.GET, "/", oas.NewOperation("test"))
		scanner.BindSchemas(openAPI)

		data, _ := json.MarshalIndent(openAPI, "", "  ")
		t.Log(string(data))
	})

	t.Run("invalid", func(t *testing.T) {
		err := tryCatch(func() {
			scanner.Def(context.Background(), pkg.TypeName("InvalidComposed"))
		})
		NewWithT(t).Expect(err).NotTo(BeNil())
	})
}

func tryCatch(fn func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("%v", e)
		}
	}()

	fn()
	return nil
}
