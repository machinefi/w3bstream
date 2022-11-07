package stringsx_test

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func Test_SplitToWords(t *testing.T) {
	words := []string{"I", "Am", "A", "10", "Years", "Senior"}
	cases := []struct {
		phrase string
		words  []string
	}{
		{"IAmA10YearsSenior", words},
		{"I Am A 10 Years Senior", words},
		{". I_ Am_A_10_Years____Senior__", words},
		{"I-~~ Am\nA\t10 Years *** Senior", words},
		{"lowercase", []string{"lowercase"}},
		{"Class", []string{"Class"}},
		{"MyClass", []string{"My", "Class"}},
		{"HTML", []string{"HTML"}},
		{"QOSType", []string{"QOS", "Type"}},
	}
	for _, c := range cases {
		NewWithT(t).Expect(SplitToWords(c.phrase)).To(Equal(c.words))
	}
}
