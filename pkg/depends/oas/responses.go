package oas

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Responses struct {
	ResponsesObject
	SpecExtensions
}

func (i Responses) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.ResponsesObject, i.SpecExtensions)
}

func (i *Responses) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.ResponsesObject, &i.SpecExtensions)
}

type ResponsesObject struct {
	Default   *Response
	Responses map[int]*Response
}

func (o *ResponsesObject) SetDefaultResponse(r *Response) {
	if r == nil {
		return
	}
	o.Default = r
}

func (o *ResponsesObject) AddResponse(statusCode int, r *Response) {
	if r == nil {
		return
	}
	if o.Responses == nil {
		o.Responses = make(map[int]*Response)
	}
	o.Responses[statusCode] = r
}

func (o ResponsesObject) MarshalJSON() ([]byte, error) {
	responses := make(map[string]*Response)
	if o.Default != nil {
		responses["default"] = o.Default
	}
	for status := range o.Responses {
		responses[fmt.Sprintf("%d", status)] = o.Responses[status]
	}
	return json.Marshal(responses)
}

func (o *ResponsesObject) UnmarshalJSON(data []byte) error {
	responses := make(map[string]*Response)
	err := json.Unmarshal(data, &responses)
	if err != nil {
		return err
	}
	for statusOrDefault := range responses {
		if statusOrDefault == "default" {
			o.Default = responses[statusOrDefault]
		} else if status, err := strconv.ParseInt(statusOrDefault, 10, 10); err == nil {
			if o.Responses == nil {
				o.Responses = make(map[int]*Response)
			}
			o.Responses[int(status)] = responses[statusOrDefault]
		}
	}
	return nil
}

func NewResponse(desc string) *Response {
	resp := &Response{}
	resp.Description = desc
	return resp
}

type Response struct {
	Reference
	ResponseObject
	SpecExtensions
}

func (r Response) MarshalJSON() ([]byte, error) {
	return r.MarshalJSONRefFirst(r.ResponseObject, r.SpecExtensions)
}

func (r *Response) UnmarshalJSON(data []byte) error {
	return r.UnmarshalJSONRefFirst(data, &r.ResponseObject, &r.SpecExtensions)
}

type ResponseObject struct {
	Description string `json:"description"`
	WithHeaders
	WithContent
	WithLinks
}

type WithLinks struct {
	Links map[string]*Link `json:"links,omitempty"`
}

func (object *WithLinks) AddLink(name string, l *Link) {
	if l == nil {
		return
	}
	if object.Links == nil {
		object.Links = make(map[string]*Link)
	}
	object.Links[name] = l
}

func NewLink(operationId string) *Link {
	return &Link{
		LinkObject: LinkObject{
			OperationId: operationId,
		},
	}
}

type Link struct {
	Reference
	LinkObject
	SpecExtensions
}

func (l Link) MarshalJSON() ([]byte, error) {
	return l.MarshalJSONRefFirst(l.LinkObject, l.SpecExtensions)
}

func (l *Link) UnmarshalJSON(data []byte) error {
	return l.UnmarshalJSONRefFirst(data, &l.LinkObject, &l.SpecExtensions)
}

type LinkObject struct {
	OperationRef string                       `json:"operationRef,omitempty"`
	OperationId  string                       `json:"operationId,omitempty"`
	Parameters   map[string]RuntimeExpression `json:"parameters,omitempty"`
	RequestBody  RuntimeExpression            `json:"requestBody,omitempty"`
	Description  string                       `json:"description,omitempty"`
	Server       *Server                      `json:"server,omitempty"`
}

func (o *LinkObject) AddParameter(name string, expr RuntimeExpression) {
	if o.Parameters == nil {
		o.Parameters = make(map[string]RuntimeExpression)
	}
	o.Parameters[name] = expr
}

func (o *LinkObject) SetRequestBody(expr RuntimeExpression) {
	o.RequestBody = expr
}

type RuntimeExpression string
