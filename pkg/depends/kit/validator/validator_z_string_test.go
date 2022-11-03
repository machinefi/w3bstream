package validator_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestString_New(t *testing.T) {
	cases := []struct {
		rule   string
		expect *String
	}{
		{
			"@string[1,1000]", &String{
				MinLength: 1,
				MaxLength: ptrx.Uint64(1000),
			},
		}, {
			"@string[1,]", &String{
				MinLength: 1,
			},
		}, {
			"@string<length>[1]", &String{
				MinLength: 1,
				MaxLength: ptrx.Uint64(1),
			},
		}, {
			"@char[1,]", &String{
				LenMode:   STR_LEN_MODE__RUNE_COUNT,
				MinLength: 1,
			},
		}, {
			"@string<rune_count>[1,]", &String{
				LenMode:   STR_LEN_MODE__RUNE_COUNT,
				MinLength: 1,
			},
		}, {
			"@string{KEY1,KEY2}", &String{
				Enums: map[string]string{
					"KEY1": "KEY1",
					"KEY2": "KEY2",
				},
			},
		}, {
			`@string/^\w+/`, &String{
				Pattern: regexp.MustCompile(`^\w+`),
			},
		}, {
			`@string/^\w+\/test/`, &String{
				Pattern: regexp.MustCompile(`^\w+/test`),
			},
		},
	}

	for i, c := range cases {
		name := fmt.Sprintf("%02d_%s|%s|%s", i, rtString, c.rule, c.expect)
		t.Run(name, func(t *testing.T) {
			v, err := c.expect.New(ctx, MustParseRuleStringByType(c.rule, rttString))
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(v).To(Equal(c.expect))
		})
	}
}

func TestString_NewFailed(t *testing.T) {
	cases := []struct {
		rule string
		rtyp typesx.Type
	}{

		{"@string", rttInt},
		{"@string<length, 1>", rttString},
		{"@string<unsupported>", rttString},
		{"@string[1,0]", rttString},
		{"@string[1,-2]", rttString},
		{"@string[a,]", rttString},
		{"@string[-1,1]", rttString},
		{"@string(-1,1)", rttString},
	}

	vldt := &String{}

	for ci, c := range cases {
		rule := MustParseRuleStringByType(c.rule, c.rtyp)
		t.Run(
			fmt.Sprintf("%02d_%s|%s", ci, c.rtyp, c.rule),
			func(t *testing.T) {
				_, err := vldt.New(ctx, rule)
				NewWithT(t).Expect(err).NotTo(BeNil())
				// t.Logf("\n%v", err)
			},
		)
	}
}

func TestString_Validate(t *testing.T) {
	cases := []struct {
		values []interface{}
		vldt   *String
		desc   string
	}{
		{
			[]interface{}{reflect.ValueOf("a"), StringType("aa"), "aaa", "aaaa", "aaaaa"}, &String{
				MaxLength: ptrx.Uint64(5),
			}, "LessThan",
		}, {
			[]interface{}{"一", "一一", "一一一"}, &String{
				LenMode:   STR_LEN_MODE__RUNE_COUNT,
				MaxLength: ptrx.Uint64(3),
			}, "CharCountLessThan",
		}, {
			[]interface{}{"A", "B"}, &String{
				Enums: map[string]string{
					"A": "A",
					"B": "B",
				},
			}, "InEnum",
		}, {
			[]interface{}{"word", "word1"}, &String{
				Pattern: regexp.MustCompile(`^\w+`),
			}, "RegexpMatched",
		},
	}

	for ci, c := range cases {
		for vi, v := range c.values {
			name := fmt.Sprintf("%02d_%02d_%s|%s|%v", ci, vi, c.desc, c.vldt, v)
			t.Run(name, func(t *testing.T) {
				err := c.vldt.Validate(v)
				NewWithT(t).Expect(err).To(BeNil())
			})
		}
	}
}

func TestString_ValidateFailed(t *testing.T) {
	cases := []struct {
		values []interface{}
		vldt   *String
		desc   string
	}{
		{
			[]interface{}{"C", "D", "E"}, &String{
				Enums: map[string]string{"A": "A", "B": "B"},
			}, "NotInEnum",
		},
		{
			[]interface{}{"-word", "-word1"}, &String{
				Pattern: regexp.MustCompile(`^\w+`),
			}, "RegexpNotMatched",
		}, {
			[]interface{}{1.1, reflect.ValueOf(1.1)}, &String{
				MinLength: 5,
			}, "UnsupportedTypes",
		}, {
			[]interface{}{"a", "aa", StringType("aaa"), []byte("aaaa")}, &String{
				MinLength: 5,
			}, "TooShort",
		}, {
			[]interface{}{"aa", "aaa", "aaaa", "aaaaa"}, &String{
				MaxLength: ptrx.Uint64(1),
			}, "TooLong",
		}, {
			[]interface{}{"字符太多"}, &String{
				LenMode:   STR_LEN_MODE__RUNE_COUNT,
				MaxLength: ptrx.Uint64(3),
			}, "TooManyChars",
		},
	}

	for ci, c := range cases {
		for vi, v := range c.values {
			name := fmt.Sprintf("%02d_%02d_%s|%s|%v", ci, vi, c.desc, c.vldt, v)
			t.Run(name, func(t *testing.T) {
				err := c.vldt.Validate(v)
				NewWithT(t).Expect(err).NotTo(BeNil())
				// t.Logf("\n%v", err)
			})
		}
	}
}
