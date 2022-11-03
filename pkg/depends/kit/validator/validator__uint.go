package validator

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type Uint struct {
	BitSize uint

	Minimum          uint64
	Maximum          uint64
	MultipleOf       uint64
	ExclusiveMaximum bool
	ExclusiveMinimum bool

	Enums map[uint64]string
}

func init() { DefaultFactory.Register(&Uint{}) }

func (Uint) Names() []string {
	return []string{"uint", "uint8", "uint16", "uint32", "uint64"}
}

func (vu *Uint) SetDefault() {
	if vu != nil {
		if vu.BitSize == 0 {
			vu.BitSize = 32
		}
		if vu.Maximum == 0 {
			vu.Maximum = MaxUint(vu.BitSize)
		}
	}
}

var TargetUintValue = "uint value"

func (vu *Uint) Validate(v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if k := rv.Type().Kind(); !typesx.IsUnsignedIntReflectKind(k) {
		return errors.NewUnsupportedTypeError(rv.Type().String(), vu.String())
	}

	val := rv.Uint()

	if vu.Enums != nil {
		if _, ok := vu.Enums[val]; !ok {
			values := make([]interface{}, 0)
			for _, v := range vu.Enums {
				values = append(values, v)
			}

			return &errors.NotInEnumError{
				Target:  TargetUintValue,
				Current: val,
				Enums:   values,
			}
		}
		return nil
	}

	if ((vu.ExclusiveMinimum && val == vu.Minimum) || val < vu.Minimum) ||
		((vu.ExclusiveMaximum && val == vu.Maximum) || val > vu.Maximum) {
		return &errors.OutOfRangeError{
			Target:           TargetUintValue,
			Current:          val,
			Minimum:          vu.Minimum,
			ExclusiveMinimum: vu.ExclusiveMinimum,
			Maximum:          vu.Maximum,
			ExclusiveMaximum: vu.ExclusiveMaximum,
		}
	}

	if vu.MultipleOf != 0 {
		if val%vu.MultipleOf != 0 {
			return &errors.MultipleOfError{
				Target:     TargetUintValue,
				Current:    val,
				MultipleOf: vu.MultipleOf,
			}
		}
	}

	return nil
}

func (Uint) New(ctx context.Context, r *Rule) (Validator, error) {
	vu := &Uint{}

	bits, err := UintRuleBitSize(r)
	if err != nil {
		return nil, err
	}
	vu.BitSize = uint(bits)

	vu.ExclusiveMinimum = r.ExclusiveMin
	vu.ExclusiveMaximum = r.ExclusiveMax

	min, max, err := UintRuleRange(r, fmt.Sprintf("uint<%d>", vu.BitSize), vu.BitSize)
	if err != nil {
		return nil, err
	}
	vu.Minimum = min
	if max != nil {
		vu.Maximum = *max
	}

	vu.SetDefault()

	multiple, enums, err := UintRuleValues(r, int(vu.BitSize))
	if err != nil {
		return nil, err
	}
	vu.MultipleOf = multiple
	vu.Enums = enums

	return vu, vu.TypeCheck(r)
}

func (vu *Uint) TypeCheck(r *Rule) error {
	switch r.Type.Kind() {
	case reflect.Uint8:
		if vu.BitSize > 8 {
			return fmt.Errorf("bit size too large for type %s", r.String())
		}
		return nil
	case reflect.Uint16:
		if vu.BitSize > 16 {
			return fmt.Errorf("bit size too large for type %s", r.String())
		}
		return nil
	case reflect.Uint, reflect.Uint32:
		if vu.BitSize > 32 {
			return fmt.Errorf("bit size too large for type %s", r.String())
		}
		return nil
	case reflect.Uint64:
		return nil
	}
	return errors.NewUnsupportedTypeError(r.String(), vu.String())
}

func (vu *Uint) String() string {
	r := rules.NewRule(vu.Names()[0])

	r.Params = []rules.Node{
		rules.NewLiteral([]byte(strconv.Itoa(int(vu.BitSize)))),
	}

	r.Range = []*rules.Lit{
		rules.NewLiteral([]byte(fmt.Sprintf("%d", vu.Minimum))),
		rules.NewLiteral([]byte(fmt.Sprintf("%d", vu.Maximum))),
	}

	r.ExclusiveMin = vu.ExclusiveMinimum
	r.ExclusiveMax = vu.ExclusiveMaximum

	if vu.MultipleOf != 0 {
		r.ValueMatrix = [][]*rules.Lit{{
			rules.NewLiteral([]byte("%" + fmt.Sprintf("%d", vu.MultipleOf))),
		}}
	} else if vu.Enums != nil {
		ruleValues := make([]*rules.Lit, 0)
		for _, e := range vu.Enums {
			ruleValues = append(ruleValues, rules.NewLiteral([]byte(e)))
		}
		r.ValueMatrix = [][]*rules.Lit{ruleValues}
	}

	return string(r.Bytes())
}
