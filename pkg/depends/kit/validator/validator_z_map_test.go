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

func TestMap_New(t *testing.T) {
	cases := []struct {
		rule   string
		expect *Map
		typ    reflect.Type
	}{{
		"@map[1,1000]", &Map{
			MinProperties: 1,
			MaxProperties: ptrx.Uint64(1000),
		}, rtMapStringString,
	}, {
		"@map<,@map[1,2]>[1,]", &Map{
			MinProperties: 1,
			ElemValidator: DefaultFactory.MustCompile(
				bg, []byte("@map[1,2]"), typesx.FromReflectType(rtMapStringString),
			),
		}, rtMapStringMapStringString,
	}, {
		"@map<@string[0,],@map[1,2]>[1,]", &Map{
			MinProperties: 1,
			KeyValidator: DefaultFactory.MustCompile(
				bg, []byte("@string[0,]"), typesx.FromReflectType(rtString),
			),
			ElemValidator: DefaultFactory.MustCompile(
				bg, []byte("@map[1,2]"), typesx.FromReflectType(rtMapStringString),
			),
		}, rtMapStringMapStringString,
	},
	}

	for i, c := range cases {
		name := fmt.Sprintf("%02d_%s|%s|%s", i+1, c.typ, c.rule, c.expect)
		t.Run(name, func(t *testing.T) {
			v, err := c.expect.New(ctx, MustParseRuleStringByType(
				c.rule, typesx.FromReflectType(c.typ)),
			)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(v).To(Equal(c.expect))
		})
	}
}

func TestMap_NewFailed(t *testing.T) {
	cases := []struct {
		name string
		rule string
		typ  reflect.Type
	}{
		{"01", "@map", rtSliceString},
		{"02", "@map<1,>", rtMapStringString},
		{"03", "@map<,2>", rtMapStringString},
		{"04", "@map<1,2,3>", rtMapStringString},
		{"05", "@map[1,0]", rtMapStringString},
		{"06", "@map[1,-2]", rtMapStringString},
		{"07", "@map[a,]", rtMapStringString},
		{"08", "@map[-1,1]", rtMapStringString},
		{"09", "@map(-1,1)", rtMapStringString},
		{"10", "@map<@unknown,>", rtMapStringString},
		{"11", "@map<,@unknown>", rtMapStringString},
		{"12", "@map<@string[0,],@unknown>", rtMapStringString},
	}
	for i, c := range cases {
		rule := MustParseRuleStringByType(c.rule, typesx.FromReflectType(c.typ))
		name := fmt.Sprintf("%02d_%s|%s", i+1, c.typ, rule.Bytes())
		t.Run(name, func(t *testing.T) {
			v := &Map{}
			_, err := v.New(ctx, rule)
			NewWithT(t).Expect(err).NotTo(BeNil())
			// t.Logf("\n%v", err)
		})
	}
}

func TestMap_Validate(t *testing.T) {
	cases := []struct {
		values    []interface{}
		validator *Map
		desc      string
	}{
		{
			[]interface{}{
				map[string]string{"1": "", "2": ""},
				map[string]string{"1": "", "2": "", "3": ""},
				map[string]string{"1": "", "2": "", "3": "", "4": ""},
			}, &Map{
				MinProperties: 2,
				MaxProperties: ptrx.Uint64(4),
			}, "InRange",
		}, {
			[]interface{}{
				reflect.ValueOf(map[string]string{"1": "", "2": ""}),
				map[string]string{"1": "", "2": "", "3": ""},
			}, &Map{
				MinProperties: 2,
				MaxProperties: ptrx.Uint64(4),
				KeyValidator: DefaultFactory.MustCompile(
					bg, []byte("@string[1,]"), rttString,
				),
				ElemValidator: DefaultFactory.MustCompile(
					bg, []byte("@string[1,]?"), rttString,
				),
			}, "KeyValueValidate",
		},
	}

	for ci, c := range cases {
		for vi, v := range c.values {
			t.Run(
				fmt.Sprintf(
					"%02d_%02d_%s|%s|%v",
					ci+1, vi+1, c.desc, c.validator, v,
				),
				func(t *testing.T) {
					NewWithT(t).Expect(c.validator.Validate(v)).To(BeNil())
				},
			)
		}
	}
}

func TestMap_ValidateFailed(t *testing.T) {
	cases := []struct {
		values    []interface{}
		validator *Map
		desc      string
	}{
		{
			[]interface{}{
				map[string]string{"1": ""},
				map[string]string{"1": "", "2": "", "3": "", "4": "", "5": ""},
				map[string]string{"1": "", "2": "", "3": "", "4": "", "5": "", "6": ""},
			}, &Map{
				MinProperties: 2,
				MaxProperties: ptrx.Uint64(4),
			}, "OutOfRange",
		}, {
			[]interface{}{
				map[string]string{"1": "", "2": ""},
				map[string]string{"1": "", "2": "", "3": ""},
			}, &Map{
				MinProperties: 2,
				MaxProperties: ptrx.Uint64(4),
				KeyValidator:  DefaultFactory.MustCompile(bg, []byte("@string[2,]"), rttString),
				ElemValidator: DefaultFactory.MustCompile(bg, []byte("@string[2,]"), rttString),
			}, "KeyElemValidateFailed",
		},
	}

	for ci, c := range cases {
		for vi, v := range c.values {
			t.Run(
				fmt.Sprintf("%02d_%02d_%s|%s|%v",
					ci+1, vi+1, c.desc, c.validator, v,
				),
				func(t *testing.T) {
					err := c.validator.Validate(v)
					NewWithT(t).Expect(err).NotTo(BeNil())
					// t.Logf("\n%v", err)
				})
		}
	}
}
