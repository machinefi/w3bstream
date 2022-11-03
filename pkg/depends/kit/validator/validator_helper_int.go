package validator

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
)

func MinInt(bits uint) int64 { return -(1 << (bits - 1)) }

func MaxInt(bits uint) int64 { return 1<<(bits-1) - 1 }

func IntRuleRange(r *Rule, bits uint) (*int64, *int64, error) {
	if r.Range == nil {
		return nil, nil, nil
	}

	typ := fmt.Sprintf("int<%d>", bits)

	parseInt := func(b []byte) (*int64, error) {
		if len(b) == 0 {
			return nil, nil
		}
		n, err := strconv.ParseInt(string(b), 10, int(bits))
		if err != nil {
			return nil, fmt.Errorf("%s value is not correct: %s", typ, err)
		}
		return &n, nil
	}

	ranges := r.Range
	switch len(ranges) {
	case 2:
		min, err := parseInt(ranges[0].Bytes())
		if err != nil {
			return nil, nil, fmt.Errorf("min %s", err)
		}
		max, err := parseInt(ranges[1].Bytes())
		if err != nil {
			return nil, nil, fmt.Errorf("max %s", err)
		}
		if min != nil && max != nil && *max < *min {
			return nil, nil, fmt.Errorf(
				"max %s value must be equal or large than min expect %d, current %d",
				typ, min, max,
			)
		}

		return min, max, nil
	case 1:
		min, err := parseInt(ranges[0].Bytes())
		if err != nil {
			return nil, nil, fmt.Errorf("min %s", err)
		}
		return min, min, nil
	}
	return nil, nil, nil
}

func IntRuleBitSize(r *Rule) (bits uint64, err error) {
	buf := &bytes.Buffer{}

	for _, char := range r.Name {
		if unicode.IsDigit(char) {
			buf.WriteRune(char)
		}
	}

	if buf.Len() == 0 && r.Params != nil {
		if len(r.Params) != 1 {
			err = fmt.Errorf("int should only 1 parameter, but got %d", len(r.Params))
			return
		}
		buf.Write(r.Params[0].Bytes())
	}

	if buf.Len() != 0 {
		s := buf.String()
		bits, err = strconv.ParseUint(s, 10, 8)
		if err != nil || bits > 64 {
			err = errors.NewSyntaxError("int parameter should be valid bit size, but got `%s`", s)
			return
		}
	}

	if bits == 0 {
		bits = 32
	}
	return
}

func IntRuleValues(r *Rule, bits int) (multiple int64, enums map[int64]string, err error) {
	values := r.ComputedValues()
	if values == nil {
		return
	}

	if len(values) == 1 {
		raw := values[0].Bytes()
		if raw[0] == '%' {
			raw = raw[1:]
			multiple, err = strconv.ParseInt(string(raw), 10, bits)
			if err != nil {
				err = errors.NewSyntaxError(
					"multipleOf should be a valid int%d value, but got `%s`",
					bits, raw,
				)
				return
			}
		}
	}

	if multiple == 0 {
		enums = map[int64]string{}
		for _, v := range values {
			s := string(v.Bytes())
			enumv, _err := strconv.ParseInt(s, 10, bits)
			if _err != nil {
				err = errors.NewSyntaxError(
					"enum should be a valid int%d value, but got `%s`", bits, v,
				)
				return
			}
			enums[enumv] = s
		}
	}
	return
}
