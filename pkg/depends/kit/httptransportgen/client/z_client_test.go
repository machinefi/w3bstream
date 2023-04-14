package client_test

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/client"
	. "github.com/onsi/gomega"
)

func TestOpenAPIGenerator(t *testing.T) {
	cwd, _ := os.Getwd()

	openAPISchema := &url.URL{Scheme: "file", Path: filepath.Join(cwd, "../testdata/server/cmd/app/openapi.json")}

	g := client.NewGenerator("demo", openAPISchema, client.OptionVendorImportByGoMod())

	g.Load()
	g.Output(filepath.Join(cwd, "../testdata/clients"))
}

func TestToColonPath(t *testing.T) {
	NewWithT(t).Expect(client.ToColonPath("/user/{userID}/tags/{tagID}")).To(Equal("/user/:userID/tags/:tagID"))
	NewWithT(t).Expect(client.ToColonPath("/user/{userID}")).To(Equal("/user/:userID"))
}

func TestGenEnumInt(t *testing.T) {
	cwd, _ := os.Getwd()
	g := client.NewGenerator("demo", &url.URL{}, client.OptionVendorImportByGoMod())
	spec := []byte(`
{
  "openapi": "3.0.3",
  "components": {
    "schemas": {
      "ExampleComCloudchainSrvDemoPkgConstantsErrorsStatusError": {
        "type": "integer",
        "format": "int32",
        "enum": [
          400000001,
          400000002
        ],
        "x-enum-labels": [
          "400000001",
          "400000002"
        ],
        "x-go-vendor-type": "example.com/cloudchain/srv-demo/pkg/constants/errors.StatusError",
        "x-id": "ExampleComCloudchainSrvDemoPkgConstantsErrorsStatusError"
      }
    }
  }
}
`)
	if err := json.NewDecoder(bytes.NewBuffer(spec)).Decode(g.Spec); err != nil {
		panic(err)
	}
	g.Output(filepath.Join(cwd, "../testdata/clients/enum_int"))
}

func TestGenEnumFloat(t *testing.T) {
	cwd, _ := os.Getwd()
	g := client.NewGenerator("demo", &url.URL{}, client.OptionVendorImportByGoMod())
	snippet := []byte(`
{
  "openapi": "3.0.3",
  "components": {
    "schemas": {
      "ExampleComCloudchainSrvDemoPkgConstantsErrorsStatusError": {
        "type": "number",
        "format": "double",
        "enum": [
          40000000.1,
          40000000.2
        ],
        "x-enum-labels": [
          "40000000.1",
          "40000000.2"
        ],
        "x-go-vendor-type": "example.com/cloudchain/srv-demo/pkg/constants/errors.StatusError",
        "x-id": "ExampleComCloudchainSrvDemoPkgConstantsErrorsStatusError"
      }
    }
  }
}
`)
	if err := json.NewDecoder(bytes.NewBuffer(snippet)).Decode(g.Spec); err != nil {
		panic(err)
	}
	g.Output(filepath.Join(cwd, "../testdata/clients/enum_float"))
}

func TestDegradation(t *testing.T) {
	cwd, _ := os.Getwd()
	g := client.NewGenerator("degradationDemo", &url.URL{}, client.OptionVendorImportByGoMod())
	snippet := []byte(`
{
  "openapi": "3.0.3",
  "info": {
    "title": "",
    "version": ""
  },
  "paths": {
    "/peer/version": {
      "get": {
        "tags": [
          "routes"
        ],
        "operationId": "DemoApi",
        "responses": {
          "200": {
            "description": "",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/DemoApiResp"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "DemoApiResp": {
        "type": "object",
        "properties": {
          "info": {
            "allOf": [
              {
                "$ref": "#/components/schemas/GithubComGolangJwtJwtV4RegisteredClaims"
              },
              {
                "x-go-field-name": "Info",
                "x-tag-json": "info"
              }
            ]
          }
        },
        "required": [
          "info"
        ],
        "x-id": "DemoApiResp"
      },
      "GithubComGolangJwtJwtV4RegisteredClaims": {
        "type": "object",
        "x-go-vendor-type": "github.com/golang-jwt/jwt/v4.RegisteredClaims",
        "x-id": "GithubComGolangJwtJwtV4RegisteredClaims"
      }
    }
  }
}
`)

	if err := json.NewDecoder(bytes.NewBuffer(snippet)).Decode(g.Spec); err != nil {
		panic(err)
	}
	g.Output(filepath.Join(cwd, "../testdata/clients"))
}
