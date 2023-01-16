package oas

import (
	"encoding/json"
)

type Paths struct {
	Paths map[string]*PathItem
	SpecExtensions
}

func (p *Paths) AddOperation(method HttpMethod, path string, op *Operation) {
	if p.Paths == nil {
		p.Paths = make(map[string]*PathItem)
	}
	if p.Paths[path] == nil {
		p.Paths[path] = &PathItem{}
	}
	p.Paths[path].AddOperation(method, op)
}

func (p Paths) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(p.Paths, p.SpecExtensions)
}

func (p *Paths) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &p.Paths, &p.SpecExtensions)
}

type PathItem struct {
	Operations
	PathItemObject
	SpecExtensions
}

func (i PathItem) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.Operations, i.PathItemObject, i.SpecExtensions)
}

func (i *PathItem) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.Operations, &i.PathItemObject, &i.SpecExtensions)
}

type HttpMethod string

const (
	GET     HttpMethod = "get"
	PUT     HttpMethod = "put"
	POST    HttpMethod = "post"
	DELETE  HttpMethod = "delete"
	OPTIONS HttpMethod = "options"
	HEAD    HttpMethod = "head"
	PATCH   HttpMethod = "patch"
	TRACE   HttpMethod = "trace"
)

type Operations struct {
	Operations map[HttpMethod]*Operation
}

func (v *Operations) AddOperation(method HttpMethod, op *Operation) {
	if v == nil {
		return
	}
	if v.Operations == nil {
		v.Operations = make(map[HttpMethod]*Operation)
	}
	v.Operations[method] = op
}

func (v Operations) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Operations)
}

func (v *Operations) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.Operations)
}

type PathItemObject struct {
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`

	Servers    []*Server    `json:"servers,omitempty"`
	Parameters []*Parameter `json:"parameters,omitempty"`
}
