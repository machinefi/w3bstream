package validator_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestFloat_New(t *testing.T) {
	cases := []struct {
		name   string
		rule   string
		typ    reflect.Type
		expect *validator.Float
	}{
		{
			"01", "@float[1,1000]", rtFloat32,
			&validator.Float{
				Minimum: ptrx.Float64(1),
				Maximum: ptrx.Float64(1000),
			},
		},
		{
			"02", "@float[1,1000]", rtFloat64,
			&validator.Float{
				Minimum: ptrx.Float64(1),
				Maximum: ptrx.Float64(1000),
			},
		},
		{
			"03", "@float32[1,1000]", rtFloat64,
			&validator.Float{
				Minimum: ptrx.Float64(1),
				Maximum: ptrx.Float64(1000),
			},
		},
		{
			"04", "@double[1,1000]", rtFloat64,
			&validator.Float{
				MaxDigits: 15,
				Minimum:   ptrx.Float64(1),
				Maximum:   ptrx.Float64(1000),
			},
		},
		{
			"05", "@float64[1,1000]", rtFloat64,
			&validator.Float{
				MaxDigits: 15,
				Minimum:   ptrx.Float64(1),
				Maximum:   ptrx.Float64(1000),
			},
		},
		{
			"06", "@float(1,1000]", rtFloat64,
			&validator.Float{
				Minimum:          ptrx.Float64(1),
				ExclusiveMinimum: true,
				Maximum:          ptrx.Float64(1000),
			},
		},
		{
			"07", "@float[.1,]", rtFloat64,
			&validator.Float{
				Minimum: ptrx.Float64(.1),
			},
		},
		{
			"08", "@float[,-1]", rtFloat64,
			&validator.Float{
				Maximum: ptrx.Float64(-1),
			},
		},
		{
			"09", "@float[-1]", rtFloat64,
			&validator.Float{
				Minimum: ptrx.Float64(-1),
				Maximum: ptrx.Float64(-1),
			},
		},
		{
			"10", "@float{1,2}", rtFloat64,
			&validator.Float{
				Enums: map[float64]string{
					1: "1",
					2: "2",
				},
			},
		},
		{
			"11", "@float{%2.2}", rtFloat64,
			&validator.Float{
				MultipleOf: 2.2,
			},
		},
		{
			"12", "@float<10,3>[1.333,2.333]", rtFloat64,
			&validator.Float{
				MaxDigits:     10,
				DecimalDigits: ptrx.Uint(3),
				Minimum:       ptrx.Float64(1.333),
				Maximum:       ptrx.Float64(2.333),
			},
		},
	}

	for _, c := range cases {
		c.expect.SetDefault()
		name := fmt.Sprintf("%s_%s%s|%s", c.name, c.typ, c.rule, c.expect.String())
		t.Run(name, func(t *testing.T) {
			ctx := validator.ContextWithFactory(
				bg, validator.DefaultFactory,
			)
			r, err := validator.ParseRuleByType(
				[]byte(c.rule),
				typesx.FromReflectType(c.typ),
			)
			NewWithT(t).Expect(err).To(BeNil())
			v, err := c.expect.New(ctx, r)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(v).To(Equal(c.expect))
		})
	}
}

func TestFloat_Failed(t *testing.T) {
	cases := []struct {
		typ  reflect.Type
		rule string
	}{
		{rtInt, `@float64`},
		{rtFloat32, `@float64`},
		{rtFloat32, `@double`},
		{rtFloat32, `@float<9>`},
		{rtFloat64, "@float<11,22,33>"},
		{rtFloat64, "@float<32,2123>"},
		{rtFloat64, "@float<@string>"},
		{rtFloat64, "@float<66>"},
		{rtFloat64, "@float<7,7>"},
		{rtFloat64, "@float[1,0]"},
		{rtFloat64, "@float[1,-2]"},
		{rtFloat64, "@float<7,2>[1.333,2]"},
		{rtFloat64, "@float<7,2>[111111.33,]"},
		{rtFloat64, "@float[a,]"},
		{rtFloat64, "@float[,a]"},
		{rtFloat64, "@float[a]"},
		{rtFloat64, `@float{%a}`},
		{rtFloat64, `@float{A,B,C}`},
	}

	for i, c := range cases {
		rule, err := validator.ParseRuleByType(
			[]byte(c.rule),
			typesx.FromReflectType(c.typ),
		)
		NewWithT(t).Expect(err).To(BeNil())

		t.Run(
			fmt.Sprintf("%02d|%s|%s", i+1, c.typ, rule.Bytes()),
			func(t *testing.T) {
				ctx := validator.ContextWithFactory(
					bg, validator.DefaultFactory,
				)
				v := &validator.Float{}
				_, err := v.New(ctx, rule)
				NewWithT(t).Expect(err).NotTo(BeNil())
				// t.Logf("\n%v", err)
			},
		)
	}
}

func TestFloat_Validate(t *testing.T) {
	cases := []struct {
		values    []interface{}
		validator *validator.Float
		desc      string
	}{
		{
			[]interface{}{reflect.ValueOf(float64(1)), float64(2), float64(3)},
			&validator.Float{
				Enums: map[float64]string{
					1: "1",
					2: "2",
					3: "3",
				},
			},
			"InEnum",
		},
		{
			[]interface{}{float64(2), float64(3), float64(4)},
			&validator.Float{
				Minimum: ptrx.Float64(2),
				Maximum: ptrx.Float64(4),
			},
			"InRange",
		},
		{
			[]interface{}{float64(2), float64(3), float64(4), float64(4)},
			&validator.Float{
				Minimum: ptrx.Float64(2),
				Maximum: ptrx.Float64(4),
			},
			"IntTypes",
		},
		{
			[]interface{}{float32(2), float32(3), float32(4)},
			&validator.Float{
				Minimum: ptrx.Float64(2),
				Maximum: ptrx.Float64(4),
			},
			"InRange",
		},
		{
			[]interface{}{-2.2, 4.4, -6.6},
			&validator.Float{
				MultipleOf: 2.2,
			},
			"MultipleOf",
		},
	}
	for ci := range cases {
		c := cases[ci]
		c.validator.SetDefault()
		for vi, v := range c.values {
			name := fmt.Sprintf(
				"%02d_%02d_%s|%s|%v",
				ci+1, vi+1, c.desc, c.validator, v,
			)
			t.Run(name, func(t *testing.T) {
				NewWithT(t).Expect(c.validator.Validate(v)).To(BeNil())
			})
		}
	}
}

func TestFloat_ValidateFailed(t *testing.T) {
	cases := []struct {
		values    []interface{}
		validator *validator.Float
		desc      string
	}{
		{[]interface{}{
			uint(2),
			"string",
			reflect.ValueOf("1"),
		}, &validator.Float{}, "unsupported type"},
		{[]interface{}{1.11, 1.22, float64(111111), float64(222221), 222.33}, &validator.Float{
			MaxDigits:     5,
			DecimalDigits: ptrx.Uint(1),
		}, "digits out out range range"},
		{[]interface{}{float64(4), float64(5), float64(6)}, &validator.Float{
			Enums: map[float64]string{
				1: "1",
				2: "2",
				3: "3",
			},
		}, "not in enum"},
		{[]interface{}{float64(1), float64(4), float64(5)}, &validator.Float{
			Minimum:          ptrx.Float64(2),
			Maximum:          ptrx.Float64(4),
			ExclusiveMaximum: true,
		}, "not in range"},
		{[]interface{}{1.1, 1.2, 1.3}, &validator.Float{
			MultipleOf: 2,
		}, "not multiple of"},
	}

	for ci, c := range cases {
		c.validator.SetDefault()
		for vi, v := range c.values {
			t.Run(
				fmt.Sprintf(
					"%02d_%02d_%s|%s|%v",
					ci+1, vi+1, c.desc, c.validator, v,
				),
				func(t *testing.T) {
					err := c.validator.Validate(v)
					NewWithT(t).Expect(err).NotTo(BeNil())
					// t.Logf("\n%v", err)
				},
			)
		}
	}
}

func TestFloat(t *testing.T) {
	floats := [][]float64{
		{99999.99999, 10, 5},
		{-0.19999999999999998, 17, 17},
		{9223372036854775808, 19, 0},
		{340282346638528859811704183484516925440, 39, 0},
		{math.MaxFloat64, 309, 0},
		{math.SmallestNonzeroFloat64, 324, 324},
	}

	for i := range floats {
		v := floats[i][0]
		n, d := validator.FloatLengthOfDigit(v)
		NewWithT(t).Expect(float64(n)).To(Equal(floats[i][1]))
		NewWithT(t).Expect(float64(d)).To(Equal(floats[i][2]))
	}
}
