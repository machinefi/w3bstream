package validator

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type Float struct {
	MaxDigits     uint
	DecimalDigits *uint

	Minimum          *float64
	Maximum          *float64
	ExclusiveMaximum bool
	ExclusiveMinimum bool

	MultipleOf float64

	Enums map[float64]string
}

func init() { DefaultFactory.Register(&Float{}) }

func (vf *Float) SetDefault() {
	if vf != nil {
		if vf.MaxDigits == 0 {
			vf.MaxDigits = 7
		}
		if vf.DecimalDigits == nil {
			vf.DecimalDigits = ptrx.Uint(2)
		}
	}
}

func (Float) Names() []string {
	return []string{"float", "double", "float32", "float64"}
}

var (
	TargetFloatValue                = "float value"
	TargetDecimalDigitsOfFloatValue = "decimal digits of float value"
	TargetTotalDigitsOfFloatValue   = "total digits of float value"
)

func (vf *Float) Validate(v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if k := rv.Type().Kind(); !typesx.IsFloatReflectKind(k) {
		return errors.NewUnsupportedTypeError(rv.Type().String(), vf.String())
	}

	val := rv.Float()
	decimal := *vf.DecimalDigits

	m, d := FloatLengthOfDigit(val)
	if m > vf.MaxDigits {
		return &errors.OutOfRangeError{
			Target:  TargetTotalDigitsOfFloatValue,
			Current: m,
			Maximum: vf.MaxDigits,
		}
	}

	if d > decimal {
		return &errors.OutOfRangeError{
			Target:  TargetDecimalDigitsOfFloatValue,
			Current: d,
			Maximum: decimal,
		}
	}

	if vf.Enums != nil {
		if _, ok := vf.Enums[val]; !ok {
			values := make([]interface{}, 0)
			for _, v := range vf.Enums {
				values = append(values, v)
			}

			return &errors.NotInEnumError{
				Target:  TargetFloatValue,
				Current: v,
				Enums:   values,
			}
		}
		return nil
	}

	if vf.Minimum != nil {
		minimum := *vf.Minimum
		if (vf.ExclusiveMinimum && val == minimum) || val < minimum {
			return &errors.OutOfRangeError{
				Target:           TargetFloatValue,
				Current:          val,
				Minimum:          minimum,
				ExclusiveMinimum: vf.ExclusiveMinimum,
			}
		}
	}

	if vf.Maximum != nil {
		maximum := *vf.Maximum
		if (vf.ExclusiveMaximum && val == maximum) || val > maximum {
			return &errors.OutOfRangeError{
				Target:           TargetFloatValue,
				Current:          val,
				Maximum:          maximum,
				ExclusiveMaximum: vf.ExclusiveMaximum,
			}
		}
	}

	if vf.MultipleOf != 0 {
		if !IsFloatMultipleOf(val, vf.MultipleOf, decimal) {
			return &errors.MultipleOfError{
				Target:     TargetFloatValue,
				Current:    val,
				MultipleOf: vf.MultipleOf,
			}
		}
	}

	return nil
}

func (Float) New(ctx context.Context, r *Rule) (Validator, error) {
	vf := &Float{}

	switch r.Name {
	case "float", "float32":
		vf.MaxDigits = 7
	case "double", "float64":
		vf.MaxDigits = 15
	}

	// @float<max_digit,max_decimal>
	// @float<5,2> 1.111(too many decimal) 12345.67 (too many digits)
	if len(r.Params) > 0 {
		digits, decimal, err := FloatRuleParam(r)
		if err != nil {
			return nil, err
		}
		vf.MaxDigits = uint(digits)
		if decimal != nil {
			vf.DecimalDigits = ptrx.Uint(uint(*decimal))
		}
	}

	vf.SetDefault()
	vf.ExclusiveMinimum = r.ExclusiveMin
	vf.ExclusiveMaximum = r.ExclusiveMax

	// @float[min,max]
	// @float(min,max]
	// @float[min,max)
	// @float(min,max)
	min, max, err := FloatRuleRange(r, vf.MaxDigits, vf.DecimalDigits)
	if err != nil {
		return nil, err
	}
	vf.Minimum, vf.Maximum = min, max

	// @float{1.1,2.2,3.3} should be one of these
	// @float{%2.2} should be multiple of 2.2
	multiple, enums, err := FloatRuleValues(r, vf.MaxDigits, vf.DecimalDigits)
	if err != nil {
		return nil, err
	}
	vf.MultipleOf = multiple
	vf.Enums = enums

	return vf, vf.TypeCheck(r)
}

func (vf *Float) TypeCheck(r *Rule) error {
	switch r.Type.Kind() {
	case reflect.Float32:
		if vf.MaxDigits > 7 {
			return fmt.Errorf("max digits too large for type %s", r)
		}
		return nil
	case reflect.Float64:
		return nil
	}
	return errors.NewUnsupportedTypeError(r.String(), vf.String())
}

func (vf *Float) String() string {
	vf.SetDefault()

	decimal := *vf.DecimalDigits
	rule := rules.NewRule(vf.Names()[0])
	rule.Params = []rules.Node{
		rules.NewLiteral([]byte(strconv.Itoa(int(vf.MaxDigits)))),
		rules.NewLiteral([]byte(strconv.Itoa(int(decimal)))),
	}

	if vf.Minimum != nil || vf.Maximum != nil {
		rule.Range = make([]*rules.Lit, 2)

		if vf.Minimum != nil {
			rule.Range[0] = rules.NewLiteral(
				[]byte(fmt.Sprintf("%."+strconv.Itoa(int(decimal))+"f", *vf.Minimum)),
			)
		}

		if vf.Maximum != nil {
			rule.Range[1] = rules.NewLiteral(
				[]byte(fmt.Sprintf("%."+strconv.Itoa(int(decimal))+"f", *vf.Maximum)),
			)
		}

		rule.ExclusiveMin = vf.ExclusiveMinimum
		rule.ExclusiveMax = vf.ExclusiveMaximum
	}

	if vf.MultipleOf != 0 {
		rule.ValueMatrix = [][]*rules.Lit{{
			rules.NewLiteral([]byte("%" + fmt.Sprintf("%."+strconv.Itoa(int(decimal))+"f", vf.MultipleOf))),
		}}
	} else if vf.Enums != nil {
		values := make([]*rules.Lit, 0)
		for _, str := range vf.Enums {
			values = append(values, rules.NewLiteral([]byte(str)))
		}
		rule.ValueMatrix = [][]*rules.Lit{values}
	}

	return string(rule.Bytes())
}
