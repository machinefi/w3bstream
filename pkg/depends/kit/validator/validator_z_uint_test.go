package validator_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestUint_New(t *testing.T) {
	cases := []struct {
		rule   string
		expect *Uint
		rtype  reflect.Type
	}{
		{
			"@uint8", &Uint{
				BitSize: 8,
				Maximum: MaxUint(8),
			}, rtUint8,
		}, {
			"@uint16", &Uint{
				BitSize: 16,
				Maximum: MaxUint(16),
			}, rtUint16,
		}, {
			"@uint8[1,]", &Uint{
				BitSize: 8,
				Minimum: 1,
				Maximum: MaxUint(8),
			}, rtUint,
		}, {
			"@uint[1,1000)", &Uint{
				Minimum:          1,
				Maximum:          1000,
				ExclusiveMaximum: true,
			}, rtUint,
		}, {
			"@uint(1,1000]", &Uint{
				Minimum:          1,
				Maximum:          1000,
				ExclusiveMinimum: true,
			}, rtUint,
		}, {
			"@uint[1,]", &Uint{
				Minimum: 1,
				Maximum: MaxUint(32),
			}, rtUint,
		},
		{
			"@uint16{1,2}", &Uint{
				BitSize: 16,
				Enums: map[uint64]string{
					1: "1",
					2: "2",
				},
			}, rtUint,
		}, {
			"@uint16{%2}", &Uint{
				BitSize:    16,
				MultipleOf: 2,
			}, rtUint,
		}, {
			"@uint<53>", &Uint{
				BitSize: 53,
				Maximum: MaxUint(53),
			}, rtUint64,
		}, {
			"@uint64", &Uint{
				BitSize: 64,
				Maximum: MaxUint(64),
			}, rtUint64,
		},
	}

	for i, c := range cases {
		c.expect.SetDefault()
		name := fmt.Sprintf("%02d_%s|%s|%s", i, c.rtype, c.rule, c.expect)
		t.Run(name, func(t *testing.T) {
			v, err := c.expect.New(
				ctx,
				MustParseRuleStringByType(c.rule, typesx.FromReflectType(c.rtype)),
			)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(v).To(Equal(c.expect))
		})
	}
}

func TestUint_ParseFailed(t *testing.T) {
	cases := []struct {
		rule string
		rtyp reflect.Type
	}{
		{"@uint16", rtFloat32},
		{"@uint16", rtUint8},
		{"@uint", rtUint16},
		{"@uint64", rtUint32},
		{"@uint<32,2123>", rtUint64},
		{"@uint<@string>", rtUint64},
		{"@uint<66>", rtUint64},
		{"@uint[1,0]", rtUint64},
		{"@uint[1,-2]", rtUint64},
		{"@uint[a,]", rtUint64},
		{"@uint[-1,1]", rtUint64},
		{"@uint(-1,1)", rtUint64},
		{`@uint8{%a}`, rtUint64},
		{`@uint8{A,B,C}`, rtUint64},
	}

	for i, c := range cases {
		rule := MustParseRuleStringByType(c.rule, typesx.FromReflectType(c.rtyp))
		name := fmt.Sprintf("%02d_%s|%s", i, c.rtyp, c.rule)
		t.Run(name, func(t *testing.T) {
			v := &Uint{}
			_, err := v.New(ctx, rule)
			NewWithT(t).Expect(err).NotTo(BeNil())
			t.Logf("\n%v", err)
		})
	}
}

func TestUint_Validate(t *testing.T) {
	cases := []struct {
		vals []interface{}
		vldt *Uint
		desc string
	}{
		{
			[]interface{}{reflect.ValueOf(uint(1)), uint(2), uint(3)}, &Uint{
				Enums: map[uint64]string{1: "1", 2: "2", 3: "3"},
			}, "InEnum",
		}, {
			[]interface{}{uint(2), uint(3), uint(4)}, &Uint{
				Minimum: 2,
				Maximum: 4,
			}, "InRange",
		}, {
			[]interface{}{uint8(2), uint16(3), uint32(4), uint64(4)}, &Uint{
				Minimum: 2,
				Maximum: 4,
			}, "UintTypes",
		}, {
			[]interface{}{uint64(2), uint64(3), uint64(4)}, &Uint{
				BitSize: 64,
				Minimum: 2,
				Maximum: 4,
			}, "InRange",
		}, {
			[]interface{}{uint(2), uint(4), uint(6)}, &Uint{
				MultipleOf: 2,
			}, "MultipleOf",
		},
	}

	for ci, c := range cases {
		c.vldt.SetDefault()
		for vi, v := range c.vals {
			name := fmt.Sprintf("%02d_%02d_%s|%s|%v", ci, vi, c.desc, c.vldt, v)
			t.Run(name, func(t *testing.T) {
				NewWithT(t).Expect(c.vldt.Validate(v)).To(BeNil())
			})
		}
	}
}

func TestUintValidator_ValidateFailed(t *testing.T) {
	cases := []struct {
		vals []interface{}
		vldt *Uint
		desc string
	}{
		{
			[]interface{}{2, "string", reflect.ValueOf(1)}, &Uint{
				BitSize: 64,
			}, "UnsupportedType",
		}, {
			[]interface{}{uint(4), uint(5), uint(6)}, &Uint{
				Enums: map[uint64]string{1: "1", 2: "2", 3: "3"},
			}, "NotInEnum",
		}, {
			[]interface{}{uint(1), uint(4), uint(5)}, &Uint{
				Minimum:          2,
				Maximum:          4,
				ExclusiveMaximum: true,
			}, "NotInRange",
		}, {
			[]interface{}{uint(1), uint(3), uint(5)}, &Uint{
				MultipleOf: 2,
			}, "NotMultipleOf",
		},
	}

	for ci, c := range cases {
		c.vldt.SetDefault()
		for vi, v := range c.vals {
			name := fmt.Sprintf("%02d_%02d_%s|%s|%v", ci, vi, c.desc, c.vldt, v)
			t.Run(name, func(t *testing.T) {
				err := c.vldt.Validate(v)
				NewWithT(t).Expect(err).NotTo(BeNil())
				t.Logf("\n%v", err)
			})
		}
	}
}
