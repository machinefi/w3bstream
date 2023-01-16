package oas_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestSecurityScheme(t *testing.T) {
	g := NewCaseGroup("SecurityScheme")

	g.It("empty", `{"type":""}`, SecurityScheme{})

	g.It("with api key", `{"type":"apiKey","name":"AccessToken","in":"header"}`, func() *SecurityScheme {
		server := NewAPIKeySecurityScheme("AccessToken", PositionHeader)
		return server
	}())

	g.It("with http basic", `{"type":"http","scheme":"basic"}`, func() *SecurityScheme {
		server := NewHTTPSecurityScheme("basic", "")
		return server
	}())

	g.It("with http bearer", `{"type":"http","scheme":"bearer","bearerFormat":"JWT"}`, func() *SecurityScheme {
		server := NewHTTPSecurityScheme("bearer", "JWT")
		return server
	}())

	g.It("with open id connect", `{"type":"openIdConnect","openIdConnectUrl":"http://xx.com"}`, func() *SecurityScheme {
		server := NewOpenIdConnectSecurityScheme("http://xx.com")
		return server
	}())

	g.It("with open id connect", `{"type":"oauth2","flows":{"implicit":{"authorizationUrl":"https://example.com/api/oauth","tokenUrl":"https://example.com/api/oauth/token_token","refreshUrl":"https://example.com/api/oauth/reflesh_token","scopes":{"read:pets":"read your pets","write:pets":"modify pets in your account"}}}}`, func() *SecurityScheme {
		server := NewOAuth2SecurityScheme(OAuthFlowsObject{
			Implicit: NewOAuthFlow(
				"https://example.com/api/oauth",
				"https://example.com/api/oauth/token_token",
				"https://example.com/api/oauth/reflesh_token",
				map[string]string{
					"write:pets": "modify pets in your account",
					"read:pets":  "read your pets",
				},
			),
		})
		return server
	}())

	g.Run(t)
}

func TestSecurityRequirement(t *testing.T) {
	g := NewCaseGroup("SecurityRequirement")

	g.It("empty", `{}`, SecurityRequirement{})

	g.It("Non-OAuth2 Security Requirement", `{"api_key":[]}`, SecurityRequirement{"api_key": []string{}})

	g.It("OAuth2 Security Requirement",
		`{"petstore_auth":["write:pets","read:pets"]}`,
		SecurityRequirement{"petstore_auth": []string{
			"write:pets",
			"read:pets",
		}})

	g.Run(t)
}
