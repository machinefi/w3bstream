package vm

import (
	"os"

	"github.com/iotexproject/w3bstream/pkg/modules/vm/common"
	"github.com/iotexproject/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

func NewInstance(path string, opts ...common.InstanceOptionSetter) (string, error) {
	code, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	i, err := newInstanceByCode(code, opts...)
	if err != nil {
		return "", err
	}
	return AddInstance(i), nil
}

func NewInstanceWithID(path string, by string, opts ...common.InstanceOptionSetter) error {
	code, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	i, err := newInstanceByCode(code, opts...)
	if err != nil {
		return err
	}

	AddInstanceByID(by, i)
	return nil
}

func newInstanceByCode(code []byte, opts ...common.InstanceOptionSetter) (wasm.Instance, error) {
	return wasmtime.NewInstanceByCode(code, opts...)
	// return wazero.NewInstanceByCode(code, opts...)
}
