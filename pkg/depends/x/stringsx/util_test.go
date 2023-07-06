package stringsx_test

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func TestGenRandomVisibleString(t *testing.T) {
	var visibleChars []byte
	for i := byte(0); i < 95; i++ {
		visibleChars = append(visibleChars, i+32)
	}
	t.Log(string(visibleChars))

	fs := []func(int, ...byte) string{
		stringsx.GenRandomVisibleStringV1,
		stringsx.GenRandomVisibleStringV2,
	}

	for _, f := range fs {
		NewWithT(t).Expect(len(f(10))).To(Equal(10))
		s := f(20, '_')
		NewWithT(t).Expect(len(s)).To(Equal(20))
		NewWithT(t).Expect(strings.Contains(s, "_")).To(BeFalse())
		s = f(100, '_', '\'', ' ')
		NewWithT(t).Expect(len(s)).To(Equal(100))
		NewWithT(t).Expect(strings.Contains(s, "_")).To(BeFalse())
		NewWithT(t).Expect(strings.Contains(s, "'")).To(BeFalse())
		NewWithT(t).Expect(strings.Contains(s, " ")).To(BeFalse())
	}
}

func BenchmarkGenRandomVisibleString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringsx.GenRandomVisibleStringV1(8)
	}
}

func BenchmarkGenRandomVisibleStringV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stringsx.GenRandomVisibleStringV2(8)
	}
}
