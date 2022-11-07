package validator

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
)

// Float helper
// TODO some should move to mics

func FloatRuleParam(r *Rule) (digits uint64, decimal *uint64, err error) {
	if len(r.Params) > 2 {
		err = fmt.Errorf(
			"float should only 1 or 2 parameter, but got %d", len(r.Params),
		)
		return
	}

	if raw := r.Params[0].Bytes(); len(raw) > 0 {
		digits, err = strconv.ParseUint(string(raw), 10, 4)
		if err != nil {
			err = errors.NewSyntaxError(
				"decimal digits should be a uint value which less than 16, but got `%s`",
				raw,
			)
			return
		}
	}

	if len(r.Params) == 1 {
		return
	}

	if raw := r.Params[1].Bytes(); len(raw) > 0 {
		_decimal := uint64(0)
		_decimal, err = strconv.ParseUint(string(raw), 10, 4)
		if err != nil || _decimal >= digits {
			err = errors.NewSyntaxError(
				"decimal digits should be a uint value which less than %d, but got `%s`",
				digits, raw)
			return
		}
		decimal = &_decimal
	}
	return
}

func FloatRuleRange(r *Rule, digits uint, decimal *uint) (*float64, *float64, error) {
	if r.Range == nil {
		return nil, nil, nil
	}
	typ := fmt.Sprintf("float<%d>", digits)
	if decimal != nil {
		typ = fmt.Sprintf("float<%d,%d>", digits, *decimal)
	}

	parseMaybeFloat := func(b []byte) (*float64, error) {
		if len(b) == 0 {
			return nil, nil
		}
		n, err := ParseFloatValue(b, digits, decimal)
		if err != nil {
			return nil, fmt.Errorf("%s value is not correct: %s", typ, err)
		}
		return &n, nil
	}

	ranges := r.Range
	switch len(r.Range) {
	case 2:
		min, err := parseMaybeFloat(ranges[0].Bytes())
		if err != nil {
			return nil, nil, fmt.Errorf("min %s", err)
		}
		max, err := parseMaybeFloat(ranges[1].Bytes())
		if err != nil {
			return nil, nil, fmt.Errorf("max %s", err)
		}
		if min != nil && max != nil && *max < *min {
			return nil, nil, fmt.Errorf(
				"max %s value must be equal or large than min value %v, current %v",
				typ, *min, *max,
			)
		}
		return min, max, nil
	case 1:
		min, err := parseMaybeFloat(ranges[0].Bytes())
		if err != nil {
			return nil, nil, fmt.Errorf("min %s", err)
		}
		return min, min, nil
	}
	return nil, nil, nil
}

func ParseFloatValue(b []byte, digits uint, decimal *uint) (float64, error) {
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return 0, err
	}

	if b[0] == '-' {
		b = b[1:]
	}

	if b[0] == '.' {
		b = append([]byte("0"), b...)
	}

	i := bytes.IndexRune(b, '.')

	decimalDigits := digits - 1
	if decimal != nil && *decimal < digits {
		decimalDigits = *decimal
	}

	m := uint(len(b) - 1)
	if uint(len(b)-1) > digits {
		return 0, fmt.Errorf(
			"max digits should be less than %d, but got %d", decimalDigits, m,
		)
	}

	if i != -1 {
		d := uint(len(b) - i - 1)
		if d > decimalDigits {
			return 0, fmt.Errorf(
				"decimal digits should be less than %d, but got %d",
				decimalDigits, d,
			)
		}
	}
	return f, nil
}

func FloatRuleValues(r *Rule, digits uint, decimal *uint) (multiple float64, enums map[float64]string, err error) {
	values := r.ComputedValues()
	if values == nil {
		return
	}

	if len(values) == 1 {
		raw := values[0].Bytes()
		if raw[0] == '%' {
			v := raw[1:]
			multiple, err = ParseFloatValue(v, digits, decimal)
			if err != nil {
				err = errors.NewSyntaxError(
					"multipleOf should be a valid float<%d> value, but got `%s`",
					digits, v,
				)
				return
			}
		}
	}

	if multiple == 0 {
		enums = map[float64]string{}
		for _, v := range values {
			b := v.Bytes()
			enumv, _err := ParseFloatValue(b, digits, decimal)
			if _err != nil {
				err = errors.NewSyntaxError(
					"enum should be a valid float<%d> value, but got `%s`",
					digits, b,
				)
				return
			}
			enums[enumv] = string(b)
		}
	}
	return
}

func IsFloatMultipleOf(v float64, div float64, decimal uint) bool {
	f := v / div
	prec := int(decimal)
	rounded := strconv.FormatFloat(f, 'f', prec, 64)
	value, _ := strconv.ParseFloat(rounded, 64)
	return value == math.Trunc(value)
}

func FloatLengthOfDigit(f float64) (uint, uint) {
	s := strconv.FormatFloat(f, 'e', -1, 64)
	var n, d int

	parts := strings.Split(s, "e")
	nd := strings.Split(parts[0], ".")
	i := nd[0]
	n = len(i)

	if len(nd) == 2 {
		d = len(nd[1])
	}

	if len(parts) == 2 {
		switch parts[1][0] {
		case '+':
			v, _ := strconv.ParseUint(parts[1][1:], 10, 64)
			n = n + int(v)
			d = d - int(v)
			if d < 0 {
				d = 0
			}
		case '-':
			v, _ := strconv.ParseUint(parts[1][1:], 10, 64)
			n = n - int(v)
			if n <= 0 {
				n = 1
			}
			d = d + int(v)
		}
	}

	if math.Abs(f) < 1.0 {
		n = 0
	}

	return uint(n + d), uint(d)
}
