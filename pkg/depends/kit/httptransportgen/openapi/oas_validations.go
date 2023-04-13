package openapi

import (
	"context"
	"go/types"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func BindSchemaValidationByValidateBytes(s *oas.Schema, typ types.Type, raw []byte) error {
	t := typesx.FromGoType(typ)

	fvldt, err := validator.DefaultFactory.Compile(
		context.Background(), raw, t,
		func(rule validator.Modifier) {
			rule.SetDefaultValue(nil)
		},
	)
	if err != nil {
		return err
	}

	if fvldt != nil {
		BindSchemaValidationByValidator(s, fvldt)
	}

	return nil
}

func BindSchemaValidationByValidator(s *oas.Schema, v validator.Validator) {
	if validatorLoader, ok := v.(*validator.Loader); ok {
		v = validatorLoader.Validator
	}
	if s == nil {
		*s = oas.Schema{}
	}

	switch vt := v.(type) {
	case *validator.Uint:
		if len(vt.Enums) > 0 {
			for _, v := range vt.Enums {
				s.Enum = append(s.Enum, v)
			}
			return
		}

		s.Minimum = ptrx.Ptr(float64(vt.Minimum))
		s.Maximum = ptrx.Ptr(float64(vt.Maximum))
		s.ExclusiveMinimum = vt.ExclusiveMinimum
		s.ExclusiveMaximum = vt.ExclusiveMaximum
		if vt.MultipleOf > 0 {
			s.MultipleOf = ptrx.Ptr(float64(vt.MultipleOf))
		}
	case *validator.Int:
		if len(vt.Enums) > 0 {
			for _, v := range vt.Enums {
				s.Enum = append(s.Enum, v)
			}
			return
		}

		if vt.Minimum != nil {
			s.Minimum = ptrx.Ptr(float64(*vt.Minimum))
		}
		if vt.Maximum != nil {
			s.Maximum = ptrx.Ptr(float64(*vt.Maximum))
		}
		s.ExclusiveMinimum = vt.ExclusiveMinimum
		s.ExclusiveMaximum = vt.ExclusiveMaximum

		if vt.MultipleOf > 0 {
			s.MultipleOf = ptrx.Ptr(float64(vt.MultipleOf))
		}
	case *validator.Float:
		if len(vt.Enums) > 0 {
			for _, v := range vt.Enums {
				s.Enum = append(s.Enum, v)
			}
			return
		}

		if vt.Minimum != nil {
			s.Minimum = ptrx.Ptr(*vt.Minimum)
		}
		if vt.Maximum != nil {
			s.Maximum = ptrx.Ptr(*vt.Maximum)
		}
		s.ExclusiveMinimum = vt.ExclusiveMinimum
		s.ExclusiveMaximum = vt.ExclusiveMaximum

		if vt.MultipleOf > 0 {
			s.MultipleOf = ptrx.Ptr(vt.MultipleOf)
		}
	case *validator.StrFmt:
		s.Type = oas.TypeString // force to type string for TextMarshaler
		s.Format = vt.Names()[0]
	case *validator.String:
		s.Type = oas.TypeString // force to type string for TextMarshaler

		if len(vt.Enums) > 0 {
			for _, v := range vt.Enums {
				s.Enum = append(s.Enum, v)
			}
			return
		}

		s.MinLength = ptrx.Ptr(vt.MinLength)
		if vt.MaxLength != nil {
			s.MaxLength = ptrx.Ptr(*vt.MaxLength)
		}
		if vt.Pattern != nil {
			s.Pattern = vt.Pattern.String()
		}
	case *validator.Slice:
		s.MinItems = ptrx.Ptr(vt.MinItems)
		if vt.MaxItems != nil {
			s.MaxItems = ptrx.Ptr(*vt.MaxItems)
		}

		if vt.ElemValidator != nil {
			BindSchemaValidationByValidator(s.Items, vt.ElemValidator)
		}
	case *validator.Map:
		s.MinProperties = ptrx.Ptr(vt.MinProperties)
		if vt.MaxProperties != nil {
			s.MaxProperties = ptrx.Ptr(*vt.MaxProperties)
		}
		if vt.ElemValidator != nil {
			BindSchemaValidationByValidator(s.AdditionalProperties.Schema, vt.ElemValidator)
		}
	}
}
