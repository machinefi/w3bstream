package wasmtime_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/runtime/wasmtime"
)

func TestNewWasmtimeVM(t *testing.T) {
	vmid := uuid.NewString()
	vm := wasmtime.NewWasmtimeVM(vmid)

	NewWithT(t).Expect(vm.ID()).To(Equal(vmid))

	wasmCode, err := os.ReadFile("../../testdata/log.wasm")
	NewWithT(t).Expect(err).To(BeNil())

	mod, err := vm.NewModule(wasmCode)
	NewWithT(t).Expect(err).To(BeNil())

	t.Log(mod.GetABINameList())

	_ = wasmtime.NewWasmtimeInstance(vm.(*wasmtime.VM), mod.(*wasmtime.Module))
}
