package wasmtime

import (
	"os"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"
)

func TestNewWasmtimeVM(t *testing.T) {
	vmid := uuid.NewString()
	vm := NewWasmtimeVM(vmid)

	NewWithT(t).Expect(vm.ID()).To(Equal(vmid))

	wasmCode, err := os.ReadFile("./testdata/crypto.wasm")
	NewWithT(t).Expect(err).To(BeNil())

	mod, err := vm.NewModule(wasmCode)
	NewWithT(t).Expect(err).To(BeNil())

	t.Log(mod.GetABINameList())
}
