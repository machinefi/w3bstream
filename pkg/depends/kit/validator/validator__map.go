package validator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
)

var (
	TargetMapLength = "map length"
)

type Map struct {
	MinProperties uint64
	MaxProperties *uint64

	KeyValidator  Validator
	ElemValidator Validator
}

func init() { DefaultFactory.Register(&Map{}) }

func (Map) Names() []string { return []string{"map"} }

func (vm *Map) Validate(v interface{}) error {
	switch rv := v.(type) {
	case reflect.Value:
		return vm.ValidateReflectValue(rv)
	default:
		return vm.ValidateReflectValue(reflect.ValueOf(v))
	}
}

func (vm *Map) ValidateReflectValue(rv reflect.Value) error {
	length := uint64(0)
	if !rv.IsNil() {
		length = uint64(rv.Len())
	}

	if length < vm.MinProperties {
		return &errors.OutOfRangeError{
			Target:  TargetMapLength,
			Current: rv.Interface(),
			Minimum: vm.MinProperties,
		}
	}

	if vm.MaxProperties != nil && length > *vm.MaxProperties {
		return &errors.OutOfRangeError{
			Target:  TargetMapLength,
			Current: rv.Interface(),
			Maximum: vm.MaxProperties,
		}
	}

	if vm.KeyValidator != nil || vm.ElemValidator != nil {
		errs := errors.NewErrorSet()
		for _, key := range rv.MapKeys() {
			vk := key.Interface()
			if vm.KeyValidator != nil {
				err := vm.KeyValidator.Validate(vk)
				if err != nil {
					errs.AddErr(err, fmt.Sprintf("%v/key", vk))
				}
			}
			if vm.ElemValidator != nil {
				err := vm.ElemValidator.Validate(rv.MapIndex(key).Interface())
				if err != nil {
					errs.AddErr(err, fmt.Sprintf("%v", vk))
				}
			}
		}
		return errs.Err()
	}

	return nil
}

func (vm *Map) New(ctx context.Context, r *Rule) (Validator, error) {
	if r.Type.Kind() != reflect.Map {
		return nil, errors.NewUnsupportedTypeError(r.String(), vm.String())
	}

	ret := &Map{}

	if r.ExclusiveMin || r.ExclusiveMax {
		return nil, errors.NewSyntaxError(
			"range mark of %s should not be `(` or `)`", ret.Names()[0],
		)
	}

	min, max, err := UintRuleRange(r, "size of map", 64)
	if err != nil {
		return nil, err
	}
	ret.MinProperties, ret.MaxProperties = min, max

	k, v, err := MapRuleParam(ctx, r)
	if err != nil {
		return nil, err
	}
	ret.KeyValidator = k
	ret.ElemValidator = v

	return ret, nil
}

func (vm *Map) String() string {
	rule := rules.NewRule(vm.Names()[0])

	if vm.KeyValidator != nil || vm.ElemValidator != nil {
		rule.Params = make([]rules.Node, 2)

		if vm.KeyValidator != nil {
			rule.Params[0] = rules.NewLiteral([]byte(vm.KeyValidator.String()))
		}

		if vm.ElemValidator != nil {
			rule.Params[1] = rules.NewLiteral([]byte(vm.ElemValidator.String()))
		}
	}

	rule.Range = RangeFromUint(vm.MinProperties, vm.MaxProperties)

	return string(rule.Bytes())
}
