package abi

import (
	"reflect"
	"sync"

	types "github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
)

var (
	abis = sync.Map{}
)

type Factory func(types.Instance) types.Context

func RegisterABI(name string, factory interface{}) {
	abis.Store(name, factory)
}

func GetABI(i types.Instance, name string) types.Context {
	if i == nil || name == "" {
		return nil
	}
	v, ok := abis.Load(name)
	if !ok {
		return nil
	}

	abiNames := i.GetModule().GetABINameList()
	for _, abiName := range abiNames {
		if name == abiName {
			f := v.(Factory)
			return f(i)
		}
	}
	return nil
}

func GetABIList(i types.Instance) []types.Context {
	if i == nil {
		return nil
	}

	factories := map[uintptr]struct{}{}
	res := make([]types.Context, 0)

	abiNames := i.GetModule().GetABINameList()
	for _, name := range abiNames {
		v, ok := abis.Load(name)
		if !ok {
			continue
		}
		ptr := reflect.ValueOf(v).Pointer()
		if _, ok := factories[ptr]; !ok {
			res = append(res, v.(Factory)(i))
			factories[ptr] = struct{}{}
		}
	}
	return res
}
