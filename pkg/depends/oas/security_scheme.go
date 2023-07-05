package oas

type WithSecuritySchemes struct {
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty"`
}

func (o *WithSecuritySchemes) AddSecurityScheme(name string, ss *SecurityScheme) {
	if ss == nil {
		return
	}
	if o.SecuritySchemes == nil {
		o.SecuritySchemes = make(map[string]*SecurityScheme)
	}
	o.SecuritySchemes[name] = ss
}

type WithSecurityRequirement struct {
	Security []*SecurityRequirement `json:"security,omitempty"`
}

func (o *WithSecurityRequirement) AddSecurityRequirement(sr *SecurityRequirement) {
	if sr == nil {
		return
	}
	o.Security = append(o.Security, sr)
}

type SecurityRequirement map[string][]string

func NewAPIKeySecurityScheme(name string, in Position) *SecurityScheme {
	return &SecurityScheme{
		SecuritySchemeObject: SecuritySchemeObject{
			Type: SecurityTypeAPIKey,
			Name: name,
			In:   in,
		},
	}
}

func NewHTTPSecurityScheme(scheme string, bearerFormat string) *SecurityScheme {
	if bearerFormat != "" {
		scheme = "bearer"
	}
	return &SecurityScheme{
		SecuritySchemeObject: SecuritySchemeObject{
			Type:         SecurityTypeHttp,
			Scheme:       scheme,
			BearerFormat: bearerFormat,
		},
	}
}

func NewOAuth2SecurityScheme(oauthFlowsObject OAuthFlowsObject) *SecurityScheme {
	return &SecurityScheme{
		SecuritySchemeObject: SecuritySchemeObject{
			Type: SecurityTypeOAuth2,
			Flows: &OAuthFlows{
				OAuthFlowsObject: oauthFlowsObject,
			},
		},
	}
}

func NewOpenIdConnectSecurityScheme(openIdConnectUrl string) *SecurityScheme {
	return &SecurityScheme{
		SecuritySchemeObject: SecuritySchemeObject{
			Type:             SecurityTypeOpenIdConnect,
			OpenIdConnectUrl: openIdConnectUrl,
		},
	}
}

type SecurityScheme struct {
	SecuritySchemeObject
	SpecExtensions
}

func (i SecurityScheme) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.SecuritySchemeObject, i.SpecExtensions)
}

func (i *SecurityScheme) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.SecuritySchemeObject, &i.SpecExtensions)
}

type SecuritySchemeObject struct {
	Type             SecurityType `json:"type"`
	Description      string       `json:"description,omitempty"`
	Name             string       `json:"name,omitempty"`
	In               Position     `json:"in,omitempty"`
	Scheme           string       `json:"scheme,omitempty"`
	BearerFormat     string       `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows  `json:"flows,omitempty"`
	OpenIdConnectUrl string       `json:"openIdConnectUrl,omitempty"`
}

type SecurityType string

const (
	SecurityTypeAPIKey        SecurityType = "apiKey"
	SecurityTypeHttp          SecurityType = "http"
	SecurityTypeOAuth2        SecurityType = "oauth2"
	SecurityTypeOpenIdConnect SecurityType = "openIdConnect"
)

type OAuthFlows struct {
	OAuthFlowsObject
	SpecExtensions
}

func (i OAuthFlows) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.OAuthFlowsObject, i.SpecExtensions)
}

func (i *OAuthFlows) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.OAuthFlowsObject, &i.SpecExtensions)
}

type OAuthFlowsObject struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

func NewOAuthFlow(authorizationURL string, tokenURL string, refreshURL string, scopes map[string]string) *OAuthFlow {
	return &OAuthFlow{
		OAuthFlowObject: OAuthFlowObject{
			AuthorizationURL: authorizationURL,
			TokenURL:         tokenURL,
			RefreshURL:       refreshURL,
			Scopes:           scopes,
		},
	}
}

type OAuthFlow struct {
	OAuthFlowObject
	SpecExtensions
}

func (i OAuthFlow) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.OAuthFlowObject, i.SpecExtensions)
}

func (i *OAuthFlow) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.OAuthFlowObject, &i.SpecExtensions)
}

type OAuthFlowObject struct {
	AuthorizationURL string            `json:"authorizationUrl"`
	TokenURL         string            `json:"tokenUrl"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}
