package oas

import (
	"encoding/json"
)

func NewSchema(tpe Type, fmt string) *Schema {
	return &Schema{
		SchemaObject: SchemaObject{
			Type:   tpe,
			Format: fmt,
		},
	}
}

func RefSchemaByRefer(refer Refer) *Schema {
	return &Schema{
		Reference: Reference{
			Refer: refer,
		},
	}
}

func RefSchema(ref string) *Schema {
	return RefSchemaByRefer(&StringRefer{
		Ref: ref,
	})
}

func Integer() *Schema {
	return NewSchema(TypeInteger, "int32")
}

func Long() *Schema {
	return NewSchema(TypeInteger, "int64")
}

func Float() *Schema {
	return NewSchema(TypeNumber, "float")
}

func Double() *Schema {
	return NewSchema(TypeNumber, "double")
}

func String() *Schema {
	return NewSchema(TypeString, "")
}

func Byte() *Schema {
	return NewSchema(TypeString, "byte")
}

func Binary() *Schema {
	return NewSchema(TypeString, "binary")
}

func Date() *Schema {
	return NewSchema(TypeString, "date")
}

func DateTime() *Schema {
	return NewSchema(TypeString, "date-time")
}

func Password() *Schema {
	return NewSchema(TypeString, "password")
}

func Boolean() *Schema {
	return NewSchema(TypeBoolean, "")
}

func ItemsOf(items *Schema) *Schema {
	return &Schema{
		SchemaObject: SchemaObject{
			Type:  TypeArray,
			Items: items,
		},
	}
}

type Props map[string]*Schema

func ObjectOf(props Props, required ...string) *Schema {
	return &Schema{
		SchemaObject: SchemaObject{
			Type:       TypeObject,
			Properties: props,
			SchemaValidation: SchemaValidation{
				Required: required,
			},
		},
	}
}

func MapOf(s *Schema) *Schema {
	return KeyValueOf(nil, s)
}

func KeyValueOf(k *Schema, s *Schema) *Schema {
	return &Schema{
		SchemaObject: SchemaObject{
			Type: TypeObject,
			AdditionalProperties: &SchemaOrBool{
				Allows: true,
				Schema: s,
			},
			PropertyNames: k,
		},
	}
}

func AllOf(schemas ...*Schema) *Schema {
	return &Schema{
		SchemaObject: SchemaObject{
			AllOf: schemas,
		},
	}
}

func AnyOf(schemas ...*Schema) *Schema {
	return &Schema{
		SchemaObject: SchemaObject{
			AnyOf: schemas,
		},
	}
}

func OneOf(schemas ...*Schema) *Schema {
	return &Schema{
		SchemaObject: SchemaObject{
			OneOf: schemas,
		},
	}
}

func Not(schema *Schema) *Schema {
	return &Schema{
		SchemaObject: SchemaObject{
			Not: schema,
		},
	}
}

type Schema struct {
	Reference
	SchemaObject
	SpecExtensions
}

func (s Schema) WithValidation(validation *SchemaValidation) *Schema {
	s.Enum = validation.Enum

	switch s.Type {
	case TypeInteger, TypeNumber:
		s.MultipleOf = validation.MultipleOf
		s.Maximum = validation.Maximum
		s.ExclusiveMaximum = validation.ExclusiveMaximum
		s.Minimum = validation.Minimum
		s.ExclusiveMinimum = validation.ExclusiveMinimum
	case TypeString:
		s.MaxLength = validation.MaxLength
		s.MinLength = validation.MinLength
		s.Pattern = validation.Pattern
	case TypeArray:
		s.MaxItems = validation.MaxItems
		s.MinItems = validation.MinItems
		s.UniqueItems = validation.UniqueItems
	case TypeObject:
		s.MaxProperties = validation.MaxProperties
		s.MinProperties = validation.MinProperties
		if len(s.Properties) > 0 {
			s.Required = validation.Required
		}
	}
	return &s
}

func (s *Schema) SetProperty(name string, propSchema *Schema, required bool) {
	if s.Type != TypeObject {
		return
	}
	if s.Properties == nil {
		s.Properties = make(map[string]*Schema)
	}
	s.Properties[name] = propSchema
	if required {
		s.Required = append(s.Required, name)
	}
}

func (s Schema) WithDesc(desc string) *Schema {
	s.Description = desc
	return &s
}

func (s Schema) WithTitle(title string) *Schema {
	s.Title = title
	return &s
}

func (s Schema) WithDiscriminator(discriminator *Discriminator) *Schema {
	s.Discriminator = discriminator
	return &s
}

func (s Schema) MarshalJSON() ([]byte, error) {
	return s.MarshalJSONRefFirst(s.SchemaObject, s.SpecExtensions)
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	return s.UnmarshalJSONRefFirst(data, &s.SchemaObject, &s.SpecExtensions)
}

type SchemaValidation struct {
	// numbers
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty"`

	// string
	MaxLength *uint64 `json:"maxLength,omitempty"`
	MinLength *uint64 `json:"minLength,omitempty"`
	Pattern   string  `json:"pattern,omitempty"`

	// array
	MaxItems    *uint64 `json:"maxItems,omitempty"`
	MinItems    *uint64 `json:"minItems,omitempty"`
	UniqueItems bool    `json:"uniqueItems,omitempty"`

	// object
	MaxProperties *uint64  `json:"maxProperties,omitempty"`
	MinProperties *uint64  `json:"minProperties,omitempty"`
	Required      []string `json:"required,omitempty"`

	// any
	Enum []interface{} `json:"enum,omitempty"`
}

type SchemaObject struct {
	Title string `json:"title,omitempty"`

	Type   Type   `json:"type,omitempty"`
	Format string `json:"format,omitempty"`

	Items                *Schema            `json:"items,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	AdditionalProperties *SchemaOrBool      `json:"additionalProperties,omitempty"`
	PropertyNames        *Schema            `json:"propertyNames,omitempty"`

	SchemaValidation

	AllOf []*Schema `json:"allOf,omitempty"`
	AnyOf []*Schema `json:"anyOf,omitempty"`
	OneOf []*Schema `json:"oneOf,omitempty"`
	Not   *Schema   `json:"not,omitempty"`

	Description string `json:"description,omitempty"`

	Default interface{} `json:"default,omitempty"`

	Nullable      bool           `json:"nullable,omitempty"`
	Discriminator *Discriminator `json:"discriminator,omitempty"`
	ReadOnly      bool           `json:"readOnly,omitempty"`
	WriteOnly     bool           `json:"writeOnly,omitempty"`
	XML           *XML           `json:"xml,omitempty"`
	ExternalDocs  *ExternalDoc   `json:"external_docs,omitempty"`
	Example       interface{}    `json:"example,omitempty"`
	Deprecated    bool           `json:"deprecated,omitempty"`
}

type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

type XML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}

type Type string

const (
	TypeInteger Type = "integer"
	TypeNumber  Type = "number"
	TypeString  Type = "string"
	TypeBoolean Type = "boolean"

	TypeArray  Type = "array"
	TypeObject Type = "object"
)

type SchemaOrBool struct {
	Allows bool
	Schema *Schema
}

func (s *SchemaOrBool) UnmarshalJSON(data []byte) error {
	s.Allows = true
	if len(data) > 0 && data[0] == '{' {
		var schema Schema
		if err := json.Unmarshal(data, &schema); err != nil {
			return err
		}
		s.Schema = &schema
	}
	return nil
}

func (s *SchemaOrBool) MarshalJSON() ([]byte, error) {
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}
	if s.Schema == nil && !s.Allows {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}
