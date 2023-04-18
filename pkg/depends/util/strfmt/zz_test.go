package strfmt_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/util/strfmt"
)

func TestProjectNameValidator(t *testing.T) {
	var cases = []*struct {
		name  string
		input string
		succ  bool
	}{
		{"TooShort", "123", false},
		{"TooLong", "123456789012345678901234567890123", false},
		{"InvalidChar", "@xadfef", false},
		{"UpperChar", "Xadfef", false},
		{"Succ1", "123123", true},
		{"Succ2", "abcedf", true},
		{"Succ3", "abc_123", true},
	}

	vldt := strfmt.ProjectNameValidator

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.succ {
				NewWithT(t).Expect(vldt.Validate(c.input)).To(BeNil())
			} else {
				NewWithT(t).Expect(vldt.Validate(c.input)).NotTo(BeNil())
			}
		})
	}
}
