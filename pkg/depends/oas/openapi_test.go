package oas_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func ExampleOpenAPI() {
	openapi := NewOpenAPI()

	openapi.Version = "1.0.0"
	openapi.Title = "Swagger Petstore"
	openapi.License = &License{
		LicenseObject: LicenseObject{
			Name: "MIT",
		},
	}

	openapi.AddTag(nil)
	openapi.AddTag(NewTag("pets"))

	openapi.AddSecurityScheme("token", NewHTTPSecurityScheme("bearer", "JWT"))

	openapi.AddServer(NewServer("http://petstore.swagger.io/v1"))

	openapi.AddSchema("Pet", ObjectOf(Props{
		"id":   Long(),
		"name": String(),
		"tag":  String(),
	}, "id", "name"))

	openapi.AddSchema("Pets", ItemsOf(openapi.RefSchema("Pet")))

	openapi.AddSchema("Error", ObjectOf(Props{
		"code":    Integer(),
		"message": String(),
	}, "code", "message"))

	{
		op := NewOperation("listPets")
		op.Summary = "List all pets"
		op.Tags = []string{"pets"}

		parameterLimit := QueryParameter("limit", Integer(), false).
			WithDesc("How many items to return at one time (max 100)")

		op.AddParameter(parameterLimit)

		{
			resp := NewResponse("An paged array of pets")

			s := String()
			s.Description = "A link to the next page of responses"
			resp.AddHeader("x-next", NewHeaderWithSchema(s))
			resp.AddContent("application/json", NewMediaTypeWithSchema(openapi.RefSchema("Pets")))

			op.AddResponse(http.StatusOK, resp)
		}

		{
			resp := NewResponse("unexpected error")
			resp.AddContent("application/json", NewMediaTypeWithSchema(openapi.RefSchema("Error")))

			op.SetDefaultResponse(resp)
		}

		openapi.AddOperation(GET, "/pets", op)
	}

	{
		op := NewOperation("createPets")
		op.Summary = "Create a pet"
		op.Tags = []string{"pets"}

		{
			resp := NewResponse("Null response")

			op.AddResponse(http.StatusNoContent, resp)
		}

		{
			resp := NewResponse("unexpected error")
			resp.AddContent("application/json", NewMediaTypeWithSchema(openapi.RefSchema("Error")))

			op.SetDefaultResponse(resp)
		}

		openapi.AddOperation(POST, "/pets", op)
	}

	data, _ := json.MarshalIndent(openapi, "\t", "\t")
	fmt.Println(string(data))
	/* Output:
	{
		"openapi": "3.0.3",
		"info": {
			"title": "Swagger Petstore",
			"license": {
				"name": "MIT"
			},
			"version": "1.0.0"
		},
		"paths": {
			"/pets": {
				"get": {
					"tags": [
						"pets"
					],
					"summary": "List all pets",
					"operationId": "listPets",
					"parameters": [
						{
							"name": "limit",
							"in": "query",
							"description": "How many items to return at one time (max 100)",
							"schema": {
								"type": "integer",
								"format": "int32"
							}
						}
					],
					"responses": {
						"200": {
							"description": "An paged array of pets",
							"headers": {
								"x-next": {
									"schema": {
										"type": "string",
										"description": "A link to the next page of responses"
									}
								}
							},
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/Pets"
									}
								}
							}
						},
						"default": {
							"description": "unexpected error",
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/Error"
									}
								}
							}
						}
					}
				},
				"post": {
					"tags": [
						"pets"
					],
					"summary": "Create a pet",
					"operationId": "createPets",
					"responses": {
						"204": {
							"description": "Null response"
						},
						"default": {
							"description": "unexpected error",
							"content": {
								"application/json": {
									"schema": {
										"$ref": "#/components/schemas/Error"
									}
								}
							}
						}
					}
				}
			}
		},
		"servers": [
			{
				"url": "http://petstore.swagger.io/v1"
			}
		],
		"tags": [
			{
				"name": "pets"
			}
		],
		"components": {
			"schemas": {
				"Error": {
					"type": "object",
					"properties": {
						"code": {
							"type": "integer",
							"format": "int32"
						},
						"message": {
							"type": "string"
						}
					},
					"required": [
						"code",
						"message"
					]
				},
				"Pet": {
					"type": "object",
					"properties": {
						"id": {
							"type": "integer",
							"format": "int64"
						},
						"name": {
							"type": "string"
						},
						"tag": {
							"type": "string"
						}
					},
					"required": [
						"id",
						"name"
					]
				},
				"Pets": {
					"type": "array",
					"items": {
						"$ref": "#/components/schemas/Pet"
					}
				}
			},
			"securitySchemes": {
				"token": {
					"type": "http",
					"scheme": "bearer",
					"bearerFormat": "JWT"
				}
			}
		}
	}
	*/
}
