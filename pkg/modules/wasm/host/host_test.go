package host_test

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/maps"

	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/proxy"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/host"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/runtime/wasmtime"
)

func TestHostFunctions(t *testing.T) {
	content, err := os.ReadFile("../testdata/data.wasm")
	NewWithT(t).Expect(err).To(BeNil())

	vmid := uuid.NewString()
	vm := wasmtime.NewWasmtimeVM(vmid)
	NewWithT(t).Expect(vm.ID()).To(Equal(vmid))

	mod, err := vm.NewModule(content)
	NewWithT(t).Expect(err).To(BeNil())

	instance := wasmtime.NewWasmtimeInstance(vm.(*wasmtime.VM), mod.(*wasmtime.Module))
	NewWithT(t).Expect(instance.ID()).To(Equal(vmid))

	prj := &models.Project{}
	prj.ProjectID = 1

	imports := types.NewDefaultImports()

	instance.SetUserdata(&proxy.ABIContext{
		Imports:  imports,
		Instance: instance,
	})

	mapping, err := host.HostFunctions(instance)
	NewWithT(t).Expect(err).To(BeNil())

	keys := maps.Keys(mapping)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		t.Logf("%s: %s", key, reflect.TypeOf(mapping[key]))
	}
}
