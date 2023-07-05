package openapi_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func TestOperatorScanner(t *testing.T) {
	cwd, _ := os.Getwd()
	pkg, _ := pkgx.LoadFrom(filepath.Join(cwd, "../testdata/router_scanner/auth"))

	scanner := openapi.NewOperatorScanner(pkg)

	cases := map[string]string{
		"RespWithDescribers": /* language=json*/ `{
  "operationId": "RespWithDescribers",
  "responses": {
    "200": {
      "description": "",
      "content": {
        "application/json": {
          "schema": {
            "type": "null"
          }
        }
      }
    },
    "499": {
      "description": "",
      "content": {
        "application/json": {
          "schema": {
            "$ref": "#/components/schemas/GithubComMachinefiW3BstreamPkgDependsKitStatusxStatusErr"
          }
        }
      },
      "x-status-errors": [
        "@StatusErr[ContextCanceled][499000000][ContextCanceled]"
      ]
    },
    "500": {
      "description": "",
      "content": {
        "application/json": {
          "schema": {
            "$ref": "#/components/schemas/GithubComMachinefiW3BstreamPkgDependsKitStatusxStatusErr"
          }
        }
      },
      "x-status-errors": [
        "@StatusErr[UnknownError][500000000][UnknownError]"
      ]
    }
  }
}`,
		"NoContent": /* language=json*/ `{
  "operationId": "NoContent",
  "responses": {
    "204": {
      "description": ""
    }
  }
}`,
		"Auth": /* language=json*/ `{
  "summary": "Auth",
  "description": "auth auth",
  "operationId": "Auth",
  "parameters": [
    {
      "name": "HBool",
      "in": "header",
      "required": true,
      "schema": {
        "type": "boolean",
        "x-go-field-name": "HBool"
      }
    },
    {
      "name": "HInt",
      "in": "header",
      "required": true,
      "schema": {
        "type": "integer",
        "format": "int32",
        "x-go-field-name": "HInt"
      }
    },
    {
      "name": "HString",
      "in": "header",
      "required": true,
      "schema": {
        "type": "string",
        "x-go-field-name": "HString"
      }
    },
    {
      "name": "bytes",
      "in": "query",
      "schema": {
        "type": "array",
        "items": {
          "type": "integer",
          "format": "uint8"
        },
        "x-go-field-name": "QBytes",
        "x-tag-name": "bytes,omitempty"
      }
    },
    {
      "name": "bytesKeep",
      "in": "query",
      "required": true,
      "schema": {
        "type": "array",
        "items": {
          "type": "integer",
          "format": "uint8"
        },
        "x-go-field-name": "QBytesKeepEmpty",
        "x-tag-name": "bytesKeep"
      }
    },
    {
      "name": "bytesOmit",
      "in": "query",
      "schema": {
        "type": "array",
        "items": {
          "type": "integer",
          "format": "uint8"
        },
        "x-go-field-name": "QBytesOmitEmpty",
        "x-tag-name": "bytesOmit,omitempty"
      }
    },
    {
      "name": "int",
      "in": "query",
      "required": true,
      "schema": {
        "type": "integer",
        "format": "int32",
        "x-go-field-name": "QInt",
        "x-tag-name": "int"
      }
    },
    {
      "name": "string",
      "in": "query",
      "required": true,
      "schema": {
        "type": "string",
        "x-go-field-name": "QString",
        "x-tag-name": "string"
      }
    },
    {
      "name": "a",
      "in": "cookie",
      "required": true,
      "schema": {
        "type": "string",
        "x-go-field-name": "CString",
        "x-tag-name": "a"
      }
    },
    {
      "name": "slice",
      "in": "cookie",
      "required": true,
      "schema": {
        "type": "array",
        "items": {
          "type": "string"
        },
        "x-go-field-name": "CSlice",
        "x-tag-name": "slice"
      }
    }
  ],
  "requestBody": {
    "required": true,
    "content": {
      "application/json": {
        "schema": {
          "allOf": [
            {
              "$ref": "#/components/schemas/Data"
            },
            {
              "x-go-field-name": "Data"
            }
          ]
        }
      }
    }
  },
  "responses": {
    "200": {
      "description": "",
      "content": {
        "*": {
          "schema": {
            "$ref": "#/components/schemas/Data"
          }
        }
      }
    }
  }
}`,
	}

	for n, result := range cases {
		t.Run(n, func(t *testing.T) {
			operation := &oas.Operation{}
			op := scanner.Operator(context.Background(), pkg.TypeName(n))
			op.BindOperation("", operation, true)
			data, _ := json.MarshalIndent(operation, "", "  ")
			t.Log(string(data))
			NewWithT(t).Expect(string(data)).To(Equal(result))
		})
	}
}
