package validator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
)

var (
	TargetStringLength = "string length"
	TargetStringValue  = "string value"
)

type String struct {
	Enums   map[string]string
	Pattern *regexp.Regexp
	LenMode StrLenMode

	MinLength uint64
	MaxLength *uint64
}

func init() { DefaultFactory.Register(&String{}) }

func (String) Names() []string { return []string{"string", "char"} }

func (vs *String) Validate(v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	_t := reflect.TypeOf("")

	if !rv.Type().ConvertibleTo(_t) {
		return errors.NewUnsupportedTypeError(rv.Type().String(), vs.String())
	}

	s := rv.Convert(_t).String()

	if vs.Enums != nil {
		if _, ok := vs.Enums[s]; !ok {
			values := make([]interface{}, 0)
			for _, v := range vs.Enums {
				values = append(values, v)
			}

			return &errors.NotInEnumError{
				Target:  TargetStringValue,
				Current: v,
				Enums:   values,
			}
		}
		return nil
	}

	if vs.Pattern != nil {
		if !vs.Pattern.MatchString(s) {
			return &errors.NotMatchError{
				Target:  TargetStringLength,
				Pattern: vs.Pattern,
				Current: v,
			}
		}
		return nil
	}

	length := StrLenModes[vs.LenMode](s)

	if length < vs.MinLength {
		return &errors.OutOfRangeError{
			Target:  TargetStringLength,
			Current: length,
			Minimum: vs.MinLength,
		}
	}

	if vs.MaxLength != nil && length > *vs.MaxLength {
		return &errors.OutOfRangeError{
			Target:  TargetStringLength,
			Current: length,
			Maximum: vs.MaxLength,
		}
	}
	return nil
}

func (String) New(ctx context.Context, r *Rule) (Validator, error) {
	vs := &String{}

	if r.ExclusiveMin || r.ExclusiveMax {
		return nil, errors.NewSyntaxError(
			"range mark of %s should not be `(` or `)`", vs.Names()[0],
		)
	}

	mode, err := StringRuleParam(r)
	if err != nil {
		return nil, err
	}
	vs.LenMode = mode

	if r.Pattern != nil {
		vs.Pattern = r.Pattern
		return vs, vs.TypeCheck(r)
	}

	min, max, err := UintRuleRange(r, fmt.Sprintf("%s of string", vs.LenMode), 64)
	if err != nil {
		return nil, err
	}
	vs.MinLength = min
	vs.MaxLength = max

	vs.Enums = StringRuleValues(r)

	return vs, vs.TypeCheck(r)
}

func (vs *String) TypeCheck(r *Rule) error {
	if r.Type.Kind() == reflect.String {
		return nil
	}
	return errors.NewUnsupportedTypeError(r.String(), vs.String())
}

func (vs *String) String() string {
	r := rules.NewRule(vs.Names()[0])

	if vs.Enums != nil {
		ruleValues := make([]*rules.Lit, 0)
		for _, e := range vs.Enums {
			ruleValues = append(ruleValues, rules.NewLiteral([]byte(e)))
		}
		r.ValueMatrix = [][]*rules.Lit{ruleValues}
	}

	r.Params = []rules.Node{
		rules.NewLiteral([]byte(vs.LenMode.String())),
	}

	if vs.Pattern != nil {
		r.Pattern = vs.Pattern
	}

	r.Range = RangeFromUint(vs.MinLength, vs.MaxLength)

	return string(r.Bytes())
}
