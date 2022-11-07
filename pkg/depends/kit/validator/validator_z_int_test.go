package validator_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestInt_New(t *testing.T) {
	cases := []struct {
		name   string
		rule   string
		typ    reflect.Type
		expect *validator.Int
	}{
		{
			"01", "@int8[1,]", rtInt8, &validator.Int{
				BitSize: 8,
				Minimum: ptrx.Int64(1),
				Maximum: ptrx.Int64(validator.MaxInt(8)),
			},
		}, {
			"02", "@int16[1,]", rtInt16, &validator.Int{
				BitSize: 16,
				Minimum: ptrx.Int64(1),
				Maximum: ptrx.Int64(validator.MaxInt(16)),
			},
		}, {
			"03", "@int[1,]", rtInt32, &validator.Int{
				Minimum: ptrx.Int64(1),
				Maximum: ptrx.Int64(validator.MaxInt(32)),
			},
		}, {
			"04", "@int[1,1000)", rtInt32, &validator.Int{
				Minimum:          ptrx.Int64(1),
				Maximum:          ptrx.Int64(1000),
				ExclusiveMaximum: true,
			},
		}, {
			"05", "@int(1,1000]", rtInt32, &validator.Int{
				Minimum:          ptrx.Int64(1),
				Maximum:          ptrx.Int64(1000),
				ExclusiveMinimum: true,
			},
		}, {
			"06", "@int[1,]", rtInt32, &validator.Int{
				Minimum: ptrx.Int64(1),
				Maximum: ptrx.Int64(validator.MaxInt(32)),
			},
		}, {
			"07", "@int[1]", rtInt32, &validator.Int{
				Minimum: ptrx.Int64(1),
				Maximum: ptrx.Int64(1),
			},
		}, {
			"08", "@int[,1]", rtInt32, &validator.Int{
				Maximum: ptrx.Int64(1),
			},
		}, {
			"09", "@int16{1,2}", rtInt32, &validator.Int{
				BitSize: 16,
				Enums:   map[int64]string{1: "1", 2: "2"},
			},
		}, {
			"10", "@int16{%2}", rtInt32, &validator.Int{
				BitSize:    16,
				MultipleOf: 2,
			},
		}, {
			"11", "@int64[1,1000]", rtInt64, &validator.Int{
				BitSize: 64,
				Minimum: ptrx.Int64(1),
				Maximum: ptrx.Int64(1000),
			},
		}, {
			"12", "@int<53>", rtInt64, &validator.Int{
				BitSize: 53,
				Maximum: ptrx.Int64(validator.MaxInt(53)),
			},
		},
	}
	for _, c := range cases {
		c.expect.SetDefault()
		name := fmt.Sprintf("%s%s%s|%s", c.name, c.typ, c.rule, c.expect.String())
		t.Run(name, func(t *testing.T) {
			ctx := validator.ContextWithFactory(bg, validator.DefaultFactory)
			r, err := validator.ParseRuleByType([]byte(c.rule), typesx.FromReflectType(c.typ))
			NewWithT(t).Expect(err).To(BeNil())
			v, err := c.expect.New(ctx, r)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(v).To(Equal(c.expect))
		})
	}
}

func TestInt_NewFailed(t *testing.T) {
	cases := []struct {
		name string
		typ  reflect.Type
		rule string
	}{
		{"01", rtFloat32, "@int16"},
		{"02", rtInt8, "@int16"},
		{"03", rtInt16, "@int"},
		{"04", rtInt32, "@int64"},
		{"05", rtInt, "@int<32,2123>"},
		{"06", rtInt, "@int<@string>"},
		{"07", rtInt, "@int<66>"},
		{"08", rtInt, "@int[1,0]"},
		{"09", rtInt, "@int[1,-2]"},
		{"10", rtInt, "@int[a,]"},
		{"11", rtInt, "@int[,a]"},
		{"12", rtInt, "@int[a]"},
		{"13", rtInt, `@int8{%a}`},
		{"14", rtInt, `@int8{A,B,C}`},
	}

	for _, c := range cases {
		r := validator.MustParseRuleStringByType(
			c.rule,
			typesx.FromReflectType(c.typ),
		)
		name := fmt.Sprintf("%s|%s|%s", c.name, c.typ, r.Bytes())
		t.Run(name, func(t *testing.T) {
			v := &validator.Int{}
			ctx := validator.ContextWithFactory(bg, validator.DefaultFactory)
			_, err := v.New(ctx, r)
			NewWithT(t).Expect(err).NotTo(BeNil())
			// t.Logf("\n%v", err)
		})
	}
}

func TestInt_Validate(t *testing.T) {
	cases := []struct {
		values    []interface{}
		validator *validator.Int
		desc      string
	}{
		{
			[]interface{}{reflect.ValueOf(int(1)), int(2), int(3)}, &validator.Int{
				Enums: map[int64]string{
					1: "1",
					2: "2",
					3: "3",
				},
			}, "InEnum",
		}, {
			[]interface{}{int(2), int(3), int(4)}, &validator.Int{
				Minimum: ptrx.Int64(2),
				Maximum: ptrx.Int64(4),
			}, "InRange",
		}, {
			[]interface{}{int8(2), int16(3), int32(4), int64(4)}, &validator.Int{
				Minimum: ptrx.Int64(2),
				Maximum: ptrx.Int64(4),
			}, "IntTypes",
		}, {
			[]interface{}{int64(2), int64(3), int64(4)}, &validator.Int{
				BitSize: 64,
				Minimum: ptrx.Int64(2),
				Maximum: ptrx.Int64(4),
			}, "InRange",
		}, {
			[]interface{}{int(2), int(4), int(6)}, &validator.Int{
				MultipleOf: 2,
			}, "MultipleOf",
		},
	}
	for ci := range cases {
		c := cases[ci]
		c.validator.SetDefault()
		for vi, v := range c.values {
			name := fmt.Sprintf("%02d%02d|%s|%s|%v",
				ci+1, vi+1, c.desc, c.validator, v)
			t.Run(name, func(t *testing.T) {
				NewWithT(t).Expect(c.validator.Validate(v)).To(BeNil())
			})
		}
	}
}

func TestInt_ValidateFailed(t *testing.T) {
	cases := []struct {
		values    []interface{}
		validator *validator.Int
		desc      string
	}{
		{
			[]interface{}{uint(2), "string", reflect.ValueOf("1")}, &validator.Int{
				BitSize: 64,
			}, "unsupported type",
		}, {
			[]interface{}{int(4), int(5), int(6)}, &validator.Int{
				Enums: map[int64]string{
					1: "1",
					2: "2",
					3: "3",
				},
			}, "not in enum",
		}, {
			[]interface{}{int(1), int(4), int(5)}, &validator.Int{
				Minimum:          ptrx.Int64(2),
				Maximum:          ptrx.Int64(4),
				ExclusiveMaximum: true,
			}, "not in range",
		}, {
			[]interface{}{int(1), int(3), int(5)}, &validator.Int{
				MultipleOf: 2,
			}, "not multiple of"},
	}

	for ci, c := range cases {
		c.validator.SetDefault()
		for vi, v := range c.values {
			name := fmt.Sprintf("%02d%02d|%s|%s|%v",
				ci+1, vi+1, c.desc, c.validator, v)
			t.Run(name, func(t *testing.T) {
				err := c.validator.Validate(v)
				NewWithT(t).Expect(err).NotTo(BeNil())
				// t.Logf("\n%v", err)
			})
		}
	}
}
