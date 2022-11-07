package validator_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestSlice_New(t *testing.T) {
	caseSet := map[reflect.Type][]struct {
		rule   string
		expect *Slice
		typ    reflect.Type
	}{
		reflect.TypeOf([]string{}): {
			{
				"@slice[1,1000]", &Slice{
					MinItems:      1,
					MaxItems:      ptrx.Uint64(1000),
					ElemValidator: DefaultFactory.MustCompile(bg, []byte(""), rttString),
				}, rtSliceString,
			}, {
				"@slice<@string[1,2]>[1,]", &Slice{
					MinItems:      1,
					ElemValidator: DefaultFactory.MustCompile(bg, []byte("@string[1,2]"), rttString),
				}, rtSliceString,
			}, {
				"@slice[1]", &Slice{
					MinItems:      1,
					MaxItems:      ptrx.Uint64(1),
					ElemValidator: DefaultFactory.MustCompile(bg, []byte(""), rttString),
				}, rtSliceString,
			},
		},
	}

	for typ, cases := range caseSet {
		for _, c := range cases {
			t.Run(fmt.Sprintf("%s %s|%s", typ, c.rule, c.expect.String()), func(t *testing.T) {
				r := MustParseRuleStringByType(c.rule, typesx.FromReflectType(typ))
				v, err := c.expect.New(ctx, r)
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(v).To(Equal(c.expect))
			})
		}
	}
}

func TestSlice_NewFailed(t *testing.T) {
	cases := []struct {
		rule string
		typ  reflect.Type
	}{
		{"@slice[2]", rtString},
		{"@slice[2]", rtArrayString},
		{"@slice<1>", rtSliceString},
		{"@slice<1,2,4>", rtSliceString},
		{"@slice[1,0]", rtSliceString},
		{"@slice[1,-2]", rtSliceString},
		{"@slice[a,]", rtSliceString},
		{"@slice[-1,1]", rtSliceString},
		{"@slice(1,1)", rtSliceString},
		{"@slice<@unknown>", rtSliceString},
		{"@array[1,2]", rtSliceString},
	}

	for i, c := range cases {
		rule := MustParseRuleStringByType(c.rule, typesx.FromReflectType(c.typ))
		t.Run(
			fmt.Sprintf("%02d_%s|%s", i+1, c.typ, rule.Bytes()),
			func(t *testing.T) {
				v := &Slice{}
				_, err := v.New(ctx, rule)
				NewWithT(t).Expect(err).NotTo(BeNil())
				// t.Logf("\n%v", err)
			},
		)
	}
}

func TestSlice_Validate(t *testing.T) {
	cases := []struct {
		values    []interface{}
		validator *Slice
		desc      string
	}{
		{
			[]interface{}{
				reflect.ValueOf([]string{"1", "2"}),
				[]string{"1", "2", "3"},
				[]string{"1", "2", "3", "4"},
			}, &Slice{
				MinItems: 2,
				MaxItems: ptrx.Uint64(4),
			}, "InRange",
		}, {
			[]interface{}{
				[]string{"1", "2"},
				[]string{"1", "2", "3"},
				[]string{"1", "2", "3", "4"},
			}, &Slice{
				MinItems: 2,
				MaxItems: ptrx.Uint64(4),
				ElemValidator: DefaultFactory.MustCompile(
					bg, []byte("@string[0,]"), rttString,
				),
			}, "ElemValidate",
		},
	}

	for ci, c := range cases {
		for vi, v := range c.values {
			name := fmt.Sprintf(
				"%02d_%02d_%s|%s|%v", ci+1, vi+1, c.desc, c.validator, v,
			)
			t.Run(name, func(t *testing.T) {
				NewWithT(t).Expect(c.validator.Validate(v)).To(BeNil())
			})
		}
	}
}

func TestSlice_ValidateFailed(t *testing.T) {
	cases := []struct {
		values    []interface{}
		validator *Slice
		desc      string
	}{
		{
			[]interface{}{
				[]string{"1"},
				[]string{"1", "2", "3", "4", "5"},
				[]string{"1", "2", "3", "4", "5", "6"},
			}, &Slice{
				MinItems: 2,
				MaxItems: ptrx.Uint64(4),
			}, "OutOfRange",
		},
		{
			[]interface{}{
				[]string{"1", "2"},
				[]string{"1", "2", "3"},
				[]string{"1", "2", "3", "4"},
			}, &Slice{
				MinItems: 2,
				MaxItems: ptrx.Uint64(4),
				ElemValidator: DefaultFactory.MustCompile(
					bg, []byte("@string[2,]"), rttString,
				),
			}, "ElemValidate",
		},
	}

	for ci, c := range cases {
		for vi, v := range c.values {
			t.Run(
				fmt.Sprintf(
					"%02d_%02d_%s|%s|%v", ci, vi, c.desc, c.validator, v,
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

func TestSlice(t *testing.T) {
	r := DefaultFactory.MustCompile(
		bg,
		[]byte("@slice<@float64<10,4>[-1.000,10000.000]?>"),
		rttSliceFloat64,
	)
	with := NewWithT(t)
	err := r.Validate([]float64{
		-0.9999,
		9999.9999,
		8.9999,
		0,
		1,
		20.1111,
	})
	with.Expect(err).To(BeNil())

	err = r.Validate([]float64{-1.1})
	with.Expect(err).NotTo(BeNil())
	// t.Logf("\n%v", err)

	err = r.Validate([]float64{10000.1})
	with.Expect(err).NotTo(BeNil())
	// t.Logf("\n%v", err)

	err = r.Validate([]float64{0.00005})
	with.Expect(err).NotTo(BeNil())
	// t.Logf("\n%v", err)
}
