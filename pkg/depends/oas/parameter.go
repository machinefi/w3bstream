package oas

func CookieParameter(name string, s *Schema, required bool) *Parameter {
	p := &Parameter{}
	p.Name = name
	p.In = PositionCookie
	p.Required = required
	p.Schema = s
	return p
}

func HeaderParameter(name string, s *Schema, required bool) *Parameter {
	p := &Parameter{}
	p.Name = name
	p.In = PositionHeader
	p.Required = required
	p.Schema = s
	return p
}

func PathParameter(name string, s *Schema) *Parameter {
	p := &Parameter{}
	p.Name = name
	p.In = PositionPath
	p.Required = true
	p.Schema = s
	return p
}

func QueryParameter(name string, s *Schema, required bool) *Parameter {
	p := &Parameter{}
	p.Name = name
	p.In = PositionQuery
	p.Required = required
	p.Schema = s
	return p
}

type WithParameters struct {
	Parameters []*Parameter `json:"parameters,omitempty"`
}

func (object *WithParameters) AddParameter(p *Parameter) {
	if p == nil {
		return
	}
	object.Parameters = append(object.Parameters, p)
}

type Parameter struct {
	Reference
	ParameterObject
	SpecExtensions
}

func (p Parameter) WithDesc(desc string) *Parameter {
	p.Description = desc
	return &p
}

func (p Parameter) MarshalJSON() ([]byte, error) {
	return p.MarshalJSONRefFirst(p.ParameterObject, p.SpecExtensions)
}

func (p *Parameter) UnmarshalJSON(data []byte) error {
	return p.UnmarshalJSONRefFirst(data, &p.ParameterObject, &p.SpecExtensions)
}

type WithHeaders struct {
	Headers map[string]*Header `json:"headers,omitempty"`
}

func (object *WithHeaders) AddHeader(name string, h *Header) {
	if h == nil {
		return
	}
	if object.Headers == nil {
		object.Headers = make(map[string]*Header)
	}
	object.Headers[name] = h
}

func NewHeaderWithSchema(s *Schema) *Header {
	h := &Header{}
	h.Schema = s
	return h
}

type Header struct {
	Reference
	ParameterCommonObject
	SpecExtensions
}

func (h Header) MarshalJSON() ([]byte, error) {
	return h.MarshalJSONRefFirst(h.ParameterCommonObject, h.SpecExtensions)
}

func (h *Header) UnmarshalJSON(data []byte) error {
	return h.UnmarshalJSONRefFirst(data, &h.ParameterCommonObject, &h.SpecExtensions)
}

type ParameterObject struct {
	Name string   `json:"name"`
	In   Position `json:"in"`
	ParameterCommonObject
}

type ParameterCommonObject struct {
	Description     string `json:"description,omitempty"`
	Required        bool   `json:"required,omitempty"`
	Deprecated      bool   `json:"deprecated,omitempty"`
	AllowEmptyValue bool   `json:"allowEmptyValue,omitempty"`

	Style         ParameterStyle `json:"style,omitempty"`
	Explode       bool           `json:"explode,omitempty"`
	AllowReserved bool           `json:"allowReserved,omitempty"`

	WithContentOrSchema
	Example interface{} `json:"example,omitempty"`
	WithExamples
}

type WithExamples struct {
	Examples map[string]*Example `json:"examples,omitempty"`
}

func (o *WithExamples) AddExample(name string, e *Example) {
	if e == nil {
		return
	}
	if o.Examples == nil {
		o.Examples = make(map[string]*Example)
	}
	o.Examples[name] = e
}

func NewExample() *Example {
	return &Example{}
}

type Example struct {
	Reference
	ExampleObject
	SpecExtensions
}

func (e Example) MarshalJSON() ([]byte, error) {
	return e.MarshalJSONRefFirst(e.ExampleObject, e.SpecExtensions)
}

func (e *Example) UnmarshalJSON(data []byte) error {
	return e.UnmarshalJSONRefFirst(data, &e.ExampleObject, &e.SpecExtensions)
}

type ExampleObject struct {
	Summary       string      `json:"summary,omitempty"`
	Description   string      `json:"description,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty"`
}

func NewRequestBody(desc string, required bool) *RequestBody {
	return &RequestBody{
		RequestBodyObject: RequestBodyObject{
			Description: desc,
			Required:    required,
		},
	}
}

type RequestBody struct {
	Reference
	RequestBodyObject
	SpecExtensions
}

func (r RequestBody) MarshalJSON() ([]byte, error) {
	return r.MarshalJSONRefFirst(r.RequestBodyObject, r.SpecExtensions)
}

func (r *RequestBody) UnmarshalJSON(data []byte) error {
	return r.UnmarshalJSONRefFirst(data, &r.RequestBodyObject, &r.SpecExtensions)
}

type RequestBodyObject struct {
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	WithContent
}

type Position string

const (
	PositionQuery  Position = "query"
	PositionPath   Position = "path"
	PositionHeader Position = "header"
	PositionCookie Position = "cookie"
)

type ParameterStyle string

const (
	// https://tools.ietf.org/html/rfc6570#section-3.2.7
	ParameterStyleMatrix ParameterStyle = "matrix"
	// https://tools.ietf.org/html/rfc6570#section-3.2.5
	ParameterStyleLabel ParameterStyle = "label"
	// https://tools.ietf.org/html/rfc6570#section-3.2.8
	ParameterStyleForm ParameterStyle = "form"
	// for array, csv https://tools.ietf.org/html/rfc6570#section-3.2.2
	ParameterStyleSimple ParameterStyle = "simple"
	// for array, ssv
	ParameterStyleSpaceDelimited ParameterStyle = "spaceDelimited"
	// for array, pipes
	ParameterStylePipeDelimited ParameterStyle = "pipeDelimited"
	// for object
	ParameterStyleDeepObject ParameterStyle = "deepObject"
)
