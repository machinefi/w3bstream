package rules_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
)

func TestSlashAndUnSlash(t *testing.T) {
	cases := []struct {
		name     string
		inAndOut [2]string
	}{
		{"1", [2]string{`/\w+\/test/`, `\w+/test`}},
		{"2", [2]string{`/a/`, `a`}},
		{"3", [2]string{`/abc/`, `abc`}},
		{"4", [2]string{`/☺/`, `☺`}},
		{"5", [2]string{`/\xFF/`, `\xFF`}},
		{"6", [2]string{`/\377/`, `\377`}},
		{"7", [2]string{`/\u1234/`, `\u1234`}},
		{"8", [2]string{`/\U00010111/`, `\U00010111`}},
		{"9", [2]string{`/\U0001011111/`, `\U0001011111`}},
		{"A", [2]string{`/\a\b\f\n\r\t\v\\\"/`, `\a\b\f\n\r\t\v\\\"`}},
		{"B", [2]string{`/\//`, `/`}},
		{"C", [2]string{`/`, ``}},
		{"D", [2]string{`/adfadf`, ``}},
	}

	for i := range cases {
		c := cases[i]
		t.Run("UnSlash:"+c.name, func(t *testing.T) {
			r, err := rules.UnSlash([]byte(c.inAndOut[0]))
			if c.inAndOut[1] == "" {
				NewWithT(t).Expect(err).NotTo(BeNil())
				return
			}
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(string(r)).To(Equal(c.inAndOut[1]))
		})
		t.Run("Slash:"+c.name, func(t *testing.T) {
			if c.inAndOut[1] == "" {
				return
			}
			NewWithT(t).Expect(string(rules.Slash([]byte(c.inAndOut[1])))).
				To(Equal(c.inAndOut[0]))
		})
	}
}

func TestParseRule(t *testing.T) {
	cases := [][]string{
		// simple
		{`@email`, `@email`},

		// with parameters
		{`@map<@email,         @url>`, `@map<@email,@url>`},
		{`@map<@string,>`, `@map<@string,>`},
		{`@map<,@string>`, `@map<,@string>`},
		{`@float32<10,6>`, `@float32<10,6>`},
		{`@float32<10,-1>`, `@float32<10,-1>`},
		{`@slice<@string>`, `@slice<@string>`},

		// with range
		{`@slice[0,   10]`, `@slice[0,10]`},
		{`@array[10]`, `@array[10]`},
		{`@string[0,)`, `@string[0,)`},
		{`@string[0,)`, `@string[0,)`},
		{`@int(0,)`, `@int(0,)`},
		{`@int(,1)`, `@int(,1)`},
		{`@float32(1.10,)`, `@float32(1.10,)`},

		// with values
		{`@string{A, B,    C}`, `@string{A,B,C}`},
		{`@string{, B,    C}`, `@string{,B,C}`},
		{`@uint{%2}`, `@uint{%2}`},

		// with value matrix
		{`@string{A, B,    C}{a,b}`, `@string{A,B,C}{a,b}`},

		// with not required mark or default value
		{`@string?`, `@string?`},
		{`@string ?`, `@string?`},
		{`@string = `, `@string = ''`},
		{`@string = '\''`, `@string = '\''`},
		{`@string = 'default value'`, `@string = 'default value'`},
		{`@string = 'defa\'ult\ value'`, `@string = 'defa\'ult\ value'`},
		{`@string = 13123`, `@string = '13123'`},
		{`@string = 1.1`, `@string = '1.1'`},

		// with regexp
		{`@string/\w+/`, `@string/\w+/`},
		{`@string/\w+     $/`, `@string/\w+     $/`},
		{`@string/\w+\/abc/`, `@string/\w+\/abc/`},
		{`@string/\w+\/\/abc/`, `@string/\w+\/\/abc/`},
		{`@string/^\w+\/test/`, `@string/^\w+\/test/`},

		// composes
		{`@string = 's'/\w+/`, `@string/\w+/ = 's'`},
		{`@map<,@string[1,]>`, `@map<,@string[1,]>`},
		{`@map<@string,>[1,2]`, `@map<@string,>[1,2]`},
		{`@map<@string = 's',>[1,2]`, `@map<@string = 's',>[1,2]`},
		{`@slice<@float64<10,4>[-1.000,100.000]?>`, `@slice<@float64<10,4>[-1.000,100.000]?>`},
	}

	for i := range cases {
		c := cases[i]
		t.Run("rule:"+c[0], func(t *testing.T) {
			r, err := rules.Parse(c[0])
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(c[1]).To(Equal(string(r.Bytes())))
		})
	}
}

func TestParseRuleFailed(t *testing.T) {
	cases := []string{
		`@`,
		`@unsupportted-name`,
		`@name<`,
		`@name[`,
		`@name(`,
		`@name{`,
		`@name/`,
		`@name)`,
		`@name<@sub[>`,
		`@name</>`,
		`@/`,
		`@name?=`,
	}

	for _, c := range cases {
		_, err := rules.Parse(c)
		NewWithT(t).Expect(err).NotTo(BeNil())
		// t.Logf("\n%v\n%v", c, err)
	}
}

//
// func TestRule(t *testing.T) {
// 	r, _ := ParseRuleString("@string{A,B,C}{a,b}{1,2}")
// 	spew.Dump(r.ComputedValues())
// }
//
