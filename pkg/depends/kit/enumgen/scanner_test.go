package enumgen_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func TestScanner(t *testing.T) {
	cwd, _ := os.Getwd()
	p, _ := pkgx.LoadFrom(filepath.Join(cwd, "./__examples__"))

	g := enumgen.New(p)

	t.Run("ScanIntStringerEnum", func(t *testing.T) {
		options, ok := g.Options(p.TypeName("Scheme"))
		NewWithT(t).Expect(ok).To(BeTrue())

		NewWithT(t).Expect(options.Len()).To(Equal(5))
		NewWithT(t).Expect(*options[0].Str).To(Equal("HTTP"))
		NewWithT(t).Expect(*options[2].Str).To(Equal("TCP"))
	})

	t.Run("ScanStringEnum", func(t *testing.T) {
		options, ok := g.Options(p.TypeName("PullPolicy"))
		NewWithT(t).Expect(ok).To(BeTrue())

		NewWithT(t).Expect(options.Len()).To(Equal(3))
		NewWithT(t).Expect(*options[0].Str).To(Equal("Always"))
		NewWithT(t).Expect(*options[2].Str).To(Equal("Never"))
	})
}
