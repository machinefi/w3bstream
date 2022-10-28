package vm

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
	w "github.com/wasmerio/wasmer-go/wasmer"
)

func TestNewWasm(t *testing.T) {
	var wasm *Wasm

	code, _ := ioutil.ReadFile("release.wasm")

	imports := []WasmImport{
		WasmImport{
			namespace: "env",
			functions: []WasmImportFunc{
				WasmImportFunc{
					name: "log",
					inputTypes: []w.ValueKind{
						w.I32,
					},
					outputTypes: []w.ValueKind{},
					nativeFunc: func(args []w.Value) ([]w.Value, error) {
						data, e := wasm.GetMemory("memory")
						if e != nil {
							return nil, e
						}
						fmt.Println(string(data[args[0].I32():]))
						return []w.Value{}, nil
					},
				},
				WasmImportFunc{
					name: "abort",
					inputTypes: []w.ValueKind{
						w.I32,
						w.I32,
						w.I32,
						w.I32,
					},
					outputTypes: []w.ValueKind{},
					nativeFunc: func(args []w.Value) ([]w.Value, error) {
						return []w.Value{}, nil
					},
				},
			},
		},
	}

	wasm, e := NewWasm(code, imports)
	NewWithT(t).Expect(e).To(BeNil())

	sum, e := wasm.ExecuteFunction("add", 1, 2)
	NewWithT(t).Expect(e).To(BeNil())

	v, ok := sum.(int32)
	NewWithT(t).Expect(ok).To(BeTrue())
	NewWithT(t).Expect(v).To(Equal(int32(3)))

	_, e = wasm.ExecuteFunction("hello")
	NewWithT(t).Expect(e).To(BeNil())
}
