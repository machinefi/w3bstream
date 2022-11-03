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

type Int struct {
	BitSize uint

	Minimum          *int64
	Maximum          *int64
	MultipleOf       int64
	ExclusiveMaximum bool
	ExclusiveMinimum bool

	Enums map[int64]string
}

func init() { DefaultFactory.Register(&Int{}) }

func (Int) Names() []string {
	return []string{"int", "int8", "int16", "int32", "int64"}
}

func (vi *Int) SetDefault() {
	if vi != nil {
		if vi.BitSize == 0 {
			vi.BitSize = 32
		}
		if vi.Maximum == nil {
			vi.Maximum = ptrx.Int64(MaxInt(vi.BitSize))
		}
		if vi.Minimum == nil {
			vi.Minimum = ptrx.Int64(MinInt(vi.BitSize))
		}
	}
}

var (
	TargetIntValue = "int value"
)

func (vi *Int) Validate(v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if k := rv.Type().Kind(); !typesx.IsSignedIntReflectKind(k) {
		return errors.NewUnsupportedTypeError(rv.Type().String(), vi.String())
	}

	val := rv.Int()

	if vi.Enums != nil {
		if _, ok := vi.Enums[val]; !ok {
			values := make([]interface{}, 0)
			for _, v := range vi.Enums {
				values = append(values, v)
			}

			return &errors.NotInEnumError{
				Target:  TargetIntValue,
				Current: val,
				Enums:   values,
			}
		}
		return nil
	}

	minimum := *vi.Minimum
	maximum := *vi.Maximum

	if ((vi.ExclusiveMinimum && val == minimum) || val < minimum) ||
		((vi.ExclusiveMaximum && val == maximum) || val > maximum) {
		return &errors.OutOfRangeError{
			Target:           TargetFloatValue,
			Current:          val,
			Minimum:          minimum,
			ExclusiveMinimum: vi.ExclusiveMinimum,
			Maximum:          maximum,
			ExclusiveMaximum: vi.ExclusiveMaximum,
		}
	}

	if vi.MultipleOf != 0 {
		if val%vi.MultipleOf != 0 {
			return &errors.MultipleOfError{
				Target:     TargetFloatValue,
				Current:    val,
				MultipleOf: vi.MultipleOf,
			}
		}
	}

	return nil
}

func (Int) New(ctx context.Context, r *Rule) (Validator, error) {
	vi := &Int{}

	bits, err := IntRuleBitSize(r)
	if err != nil {
		return nil, err
	}
	vi.BitSize = uint(bits)

	min, max, err := IntRuleRange(r, vi.BitSize)
	if err != nil {
		return nil, err
	}
	vi.Maximum, vi.Minimum = max, min

	vi.SetDefault()
	vi.ExclusiveMinimum = r.ExclusiveMin
	vi.ExclusiveMaximum = r.ExclusiveMax

	multiple, enums, err := IntRuleValues(r, int(vi.BitSize))
	if err != nil {
		return nil, err
	}
	vi.MultipleOf, vi.Enums = multiple, enums

	return vi, vi.TypeCheck(r)
}

func (vi *Int) TypeCheck(r *Rule) error {
	switch r.Type.Kind() {
	case reflect.Int8:
		if vi.BitSize > 8 {
			return fmt.Errorf("bit size too large for type %s", r.Type)
		}
		return nil
	case reflect.Int16:
		if vi.BitSize > 16 {
			return fmt.Errorf("bit size too large for type %s", r.Type)
		}
		return nil
	case reflect.Int, reflect.Int32:
		if vi.BitSize > 32 {
			return fmt.Errorf("bit size too large for type %s", r.Type)
		}
		return nil
	case reflect.Int64:
		return nil
	}
	return errors.NewUnsupportedTypeError(r.String(), vi.String())
}

func (vi *Int) String() string {
	r := rules.NewRule(vi.Names()[0])

	r.Params = []rules.Node{
		rules.NewLiteral([]byte(strconv.Itoa(int(vi.BitSize)))),
	}

	if vi.Minimum != nil || vi.Maximum != nil {
		r.Range = make([]*rules.Lit, 2)

		if vi.Minimum != nil {
			r.Range[0] = rules.NewLiteral(
				[]byte(fmt.Sprintf("%d", *vi.Minimum)),
			)
		}

		if vi.Maximum != nil {
			r.Range[1] = rules.NewLiteral(
				[]byte(fmt.Sprintf("%d", *vi.Maximum)),
			)
		}

		r.ExclusiveMin = vi.ExclusiveMinimum
		r.ExclusiveMax = vi.ExclusiveMaximum
	}

	r.ExclusiveMin = vi.ExclusiveMinimum
	r.ExclusiveMax = vi.ExclusiveMaximum

	if vi.MultipleOf != 0 {
		r.ValueMatrix = [][]*rules.Lit{{
			rules.NewLiteral([]byte("%" + fmt.Sprintf("%d", vi.MultipleOf))),
		}}
	} else if vi.Enums != nil {
		ruleValues := make([]*rules.Lit, 0)
		for _, str := range vi.Enums {
			ruleValues = append(ruleValues, rules.NewLiteral([]byte(str)))
		}
		r.ValueMatrix = [][]*rules.Lit{ruleValues}
	}

	return string(r.Bytes())
}
