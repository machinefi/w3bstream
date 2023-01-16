package oas

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Components struct {
	ComponentsObject
	SpecExtensions
}

func (i Components) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.ComponentsObject, i.SpecExtensions)
}

func (i *Components) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.ComponentsObject, &i.SpecExtensions)
}

type ComponentsObject struct {
	Schemas    map[string]*Schema    `json:"schemas,omitempty"`
	Responses  map[string]*Response  `json:"responses,omitempty"`
	Parameters map[string]*Parameter `json:"parameters,omitempty"`
	WithExamples
	RequestBodies map[string]*RequestBody `json:"requestBodies,omitempty"`
	WithHeaders
	WithSecuritySchemes
	WithLinks
	WithCallbacks
}

func (object *ComponentsObject) AddSchema(id string, s *Schema) {
	if s == nil {
		return
	}
	if object.Schemas == nil {
		object.Schemas = make(map[string]*Schema)
	}
	object.Schemas[id] = s
}

func (object *ComponentsObject) AddResponse(id string, r *Response) {
	if r == nil {
		return
	}
	if object.Responses == nil {
		object.Responses = make(map[string]*Response)
	}
	object.Responses[id] = r
}

func (object *ComponentsObject) AddParameter(id string, p *Parameter) {
	if p == nil {
		return
	}
	if object.Parameters == nil {
		object.Parameters = make(map[string]*Parameter)
	}
	object.Parameters[id] = p
}

func (object *ComponentsObject) AddRequestBody(id string, e *RequestBody) {
	if e == nil {
		return
	}
	if object.RequestBodies == nil {
		object.RequestBodies = make(map[string]*RequestBody)
	}
	object.RequestBodies[id] = e
}

func (object *ComponentsObject) RefSchema(id string) *Schema {
	if object.Schemas == nil || object.Schemas[id] == nil {
		return nil
	}
	return RefSchemaByRefer(NewComponentRefer("schemas", id))
}

func (object *ComponentsObject) RefResponse(id string) *Response {
	if object.Responses == nil || object.Responses[id] == nil {
		return nil
	}
	s := &Response{}
	s.Refer = NewComponentRefer("responses", id)
	return s
}

func (object *ComponentsObject) RefParameter(id string) *Parameter {
	if object.Parameters == nil || object.Parameters[id] == nil {
		return nil
	}
	s := &Parameter{}
	s.Refer = NewComponentRefer("parameters", id)
	return s
}

func (object *ComponentsObject) RefExample(id string) *Example {
	if object.Examples == nil || object.Examples[id] == nil {
		return nil
	}
	s := &Example{}
	s.Refer = NewComponentRefer("examples", id)
	return s
}

func (object *ComponentsObject) RefRequestBody(id string) *RequestBody {
	if object.RequestBodies == nil || object.RequestBodies[id] == nil {
		return nil
	}
	s := &RequestBody{}
	s.Refer = NewComponentRefer("requestBodies", id)
	return s
}

func (object *ComponentsObject) RefHeader(id string) *Header {
	if object.Headers == nil || object.Headers[id] == nil {
		return nil
	}
	s := &Header{}
	s.Refer = NewComponentRefer("headers", id)
	return s
}

func (object *ComponentsObject) RefLink(id string) *Link {
	if object.Links == nil || object.Headers[id] == nil {
		return nil
	}
	s := &Link{}
	s.Refer = NewComponentRefer("links", id)
	return s
}

func (object *ComponentsObject) RefCallback(id string) *Callback {
	if object.Callbacks == nil || object.Callbacks[id] == nil {
		return nil
	}
	s := &Callback{}
	s.Refer = NewComponentRefer("callbacks", id)
	return s
}

func (object *ComponentsObject) RequireSecurity(id string, scopes ...string) SecurityRequirement {
	if object.SecuritySchemes == nil || object.SecuritySchemes[id] == nil {
		return nil
	}
	ss := object.SecuritySchemes[id]
	if ss.Type == SecurityTypeOAuth2 {
		return SecurityRequirement{
			id: scopes,
		}
	}
	return SecurityRequirement{
		id: []string{},
	}
}

type Reference struct {
	Refer Refer
}

type Refer interface {
	RefString() string
}

func NewComponentRefer(group string, id string) *ComponentRefer {
	return &ComponentRefer{
		Group: group,
		ID:    id,
	}
}

func ParseComponentRefer(ref string) *ComponentRefer {
	if strings.HasPrefix(ref, "#/components") {
		parts := strings.Split(ref, "/")
		if len(parts) == 4 {
			return &ComponentRefer{
				Group: parts[2],
				ID:    parts[3],
			}
		}
	}
	return nil
}

type ComponentRefer struct {
	Group string
	ID    string
}

func (ref ComponentRefer) RefString() string {
	return fmt.Sprintf("#/components/%s/%s", ref.Group, ref.ID)
}

type StringRefer struct {
	Ref string `json:"$ref,omitempty"`
}

func (ref StringRefer) RefString() string {
	return ref.Ref
}

func (ref Reference) MarshalJSONRefFirst(values ...interface{}) ([]byte, error) {
	if ref.Refer != nil {
		return json.Marshal(&StringRefer{
			Ref: ref.Refer.RefString(),
		})
	}
	return flattenMarshalJSON(values...)
}

func (ref *Reference) UnmarshalJSONRefFirst(data []byte, values ...interface{}) error {
	r := &StringRefer{}
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	if r.Ref != "" {
		componentRefer := ParseComponentRefer(r.Ref)
		if componentRefer != nil {
			ref.Refer = componentRefer
			return nil
		}
		ref.Refer = r
		return nil
	}
	return flattenUnmarshalJSON(data, values...)
}
