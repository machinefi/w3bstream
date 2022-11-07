package validator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
)

type Slice struct {
	ElemValidator Validator

	MinItems uint64
	MaxItems *uint64
}

func init() { DefaultFactory.Register(&Slice{}) }

func (Slice) Names() []string { return []string{"slice", "array"} }

func (vs *Slice) Validate(v interface{}) error {
	switch rv := v.(type) {
	case reflect.Value:
		return vs.ValidateReflectValue(rv)
	default:
		return vs.ValidateReflectValue(reflect.ValueOf(v))
	}
}

var TargetSliceLength = "slice length"

func (vs *Slice) ValidateReflectValue(rv reflect.Value) error {
	length := uint64(0)
	if !rv.IsNil() {
		length = uint64(rv.Len())
	}
	if length < vs.MinItems {
		return &errors.OutOfRangeError{
			Target:  TargetSliceLength,
			Current: length,
			Minimum: vs.MinItems,
		}
	}
	if vs.MaxItems != nil && length > *vs.MaxItems {
		return &errors.OutOfRangeError{
			Target:  TargetSliceLength,
			Current: length,
			Maximum: vs.MaxItems,
		}
	}

	if vs.ElemValidator != nil {
		errs := errors.NewErrorSet()
		for i := 0; i < rv.Len(); i++ {
			err := vs.ElemValidator.Validate(rv.Index(i))
			if err != nil {
				errs.AddErr(err, i)
			}
		}
		return errs.Err()
	}
	return nil
}

func (Slice) New(ctx context.Context, r *Rule) (Validator, error) {
	vs := &Slice{}

	if r.ExclusiveMin || r.ExclusiveMax {
		return nil, errors.NewSyntaxError(
			"range mark of %s should not be `(` or `)`", vs.Names()[0],
		)
	}

	min, max, err := UintRuleRange(r, "length of slice", 64)
	if err != nil {
		return nil, err
	}
	vs.MinItems = min
	vs.MaxItems = max

	switch r.Type.Kind() {
	case reflect.Array:
		if r.Type.Len() != int(vs.MinItems) {
			return nil, fmt.Errorf(
				"length(%d) or rule should equal length(%d)",
				vs.MinItems, r.Type.Len(),
			)
		}
	case reflect.Slice:
	default:
		return nil, errors.NewUnsupportedTypeError(r.String(), vs.String())
	}

	elemValidator, err := SliceElementRule(ctx, r)
	if err != nil {
		return nil, err
	}
	vs.ElemValidator = elemValidator

	return vs, vs.TypeCheck(r)
}

func (vs *Slice) TypeCheck(r *Rule) error {
	switch r.Type.Kind() {
	case reflect.Array:
		if r.Type.Len() != int(vs.MinItems) {
			return fmt.Errorf(
				"length(%d) or r should equal length(%d) of array",
				vs.MinItems, r.Type.Len(),
			)
		}
	case reflect.Slice:
	default:
		return errors.NewUnsupportedTypeError(r.String(), vs.String())
	}
	return nil
}

func (vs *Slice) String() string {
	rule := rules.NewRule(vs.Names()[0])

	if vs.ElemValidator != nil {
		rule.Params = append(
			rule.Params,
			rules.NewLiteral([]byte(vs.ElemValidator.String())),
		)
	}

	rule.Range = RangeFromUint(vs.MinItems, vs.MaxItems)

	return string(rule.Bytes())
}
