package wasmtime

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

func TestDebugParseDwarf(t *testing.T) {
	content, err := os.ReadFile("../../testdata/data.wasm")
	NewWithT(t).Expect(err).To(BeNil())

	debug := ParseDwarf(content)
	NewWithT(t).Expect(debug).NotTo(BeNil())
	NewWithT(t).Expect(debug.data).NotTo(BeNil())
	NewWithT(t).Expect(debug.codeSectionOffset).To(Equal(0x326)) // code section start addr

	lr := debug.GetLineReader()
	NewWithT(t).Expect(lr).NotTo(BeNil())

	line := debug.SeekPC(uint64(0x2ef1)) // f3
	NewWithT(t).Expect(line).NotTo(BeNil())
}
