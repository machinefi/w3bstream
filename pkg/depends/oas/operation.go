package oas

func NewOperation(operationId string) *Operation {
	op := &Operation{}
	op.OperationId = operationId
	return op
}

type Operation struct {
	OperationObject
	SpecExtensions
}

func (op Operation) WithTags(tags ...string) *Operation {
	op.Tags = append(op.Tags, tags...)
	return &op
}

func (op Operation) WithSummary(summary string) *Operation {
	op.Summary = summary
	return &op
}

func (op Operation) WithDesc(desc string) *Operation {
	op.Description = desc
	return &op
}

func (op Operation) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(op.OperationObject, op.SpecExtensions)
}

func (op *Operation) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &op.OperationObject, &op.SpecExtensions)
}

type OperationObject struct {
	Tags         []string     `json:"tags,omitempty"`
	Summary      string       `json:"summary,omitempty"`
	Description  string       `json:"description,omitempty"`
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`

	OperationId string `json:"operationId"`
	WithParameters
	RequestBody *RequestBody `json:"requestBody,omitempty"`
	Responses   Responses    `json:"responses"`
	WithCallbacks
	WithSecurityRequirement
	Deprecated bool `json:"deprecated,omitempty"`
	WithServers
}

func (o *OperationObject) SetRequestBody(rb *RequestBody) {
	o.RequestBody = rb
}

func (o *OperationObject) AddResponse(statusCode int, r *Response) {
	o.Responses.AddResponse(statusCode, r)
}

func (o *OperationObject) SetDefaultResponse(r *Response) {
	o.Responses.SetDefaultResponse(r)
}

type WithCallbacks struct {
	Callbacks map[string]*Callback `json:"callbacks,omitempty"`
}

func (o *WithCallbacks) AddCallback(name string, c *Callback) {
	if c == nil {
		return
	}
	if o.Callbacks == nil {
		o.Callbacks = make(map[string]*Callback)
	}
	o.Callbacks[name] = c
}

func NewCallback(method HttpMethod, rule RuntimeExpression, op *Operation) *Callback {
	return &Callback{
		CallbackObject: CallbackObject{
			rule: &PathItem{
				Operations: Operations{
					Operations: map[HttpMethod]*Operation{
						method: op,
					},
				},
			},
		},
	}
}

type Callback struct {
	Reference
	CallbackObject
	SpecExtensions
}

func (i Callback) MarshalJSON() ([]byte, error) {
	return i.MarshalJSONRefFirst(i.CallbackObject, i.SpecExtensions)
}

func (i *Callback) UnmarshalJSON(data []byte) error {
	return i.UnmarshalJSONRefFirst(data, &i.CallbackObject, &i.SpecExtensions)
}

type CallbackObject map[RuntimeExpression]*PathItem
