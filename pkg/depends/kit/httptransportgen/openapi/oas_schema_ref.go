package openapi

import "github.com/machinefi/w3bstream/pkg/depends/oas"

func NewSchemaRefer(s *oas.Schema) *SchemaRefer {
	return &SchemaRefer{
		Schema: s,
	}
}

type SchemaRefer struct {
	*oas.Schema
}

func (r SchemaRefer) RefString() string {
	s := r.Schema
	if r.Schema.AllOf != nil {
		s = r.AllOf[len(r.Schema.AllOf)-1]
	}
	return oas.NewComponentRefer("schemas", s.Extensions[XID].(string)).RefString()
}
