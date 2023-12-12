package wasmtime

import (
	_ "embed"
	"testing"

	. "github.com/onsi/gomega"
)

var (
	//go:embed testdata/data.wasm
	content []byte
)

func TestDebugParseDwarf(t *testing.T) {
	debug := ParseDwarf(content)
	NewWithT(t).Expect(debug).NotTo(BeNil())
	NewWithT(t).Expect(debug.data).NotTo(BeNil())
	NewWithT(t).Expect(debug.codeSectionOffset).To(Equal(0x326)) // code section start addr

	lr := debug.GetLineReader()
	NewWithT(t).Expect(lr).NotTo(BeNil())

	line := debug.SeekPC(uint64(0x2ef1)) // f3
	NewWithT(t).Expect(line).NotTo(BeNil())
}
