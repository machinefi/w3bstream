package host

import (
	"golang.org/x/text/encoding/unicode"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/consts"
)

func GetImportsHandler(i types.Instance) types.ImportsHandler {
	ctx := GetContext(i)
	if ctx == nil {
		return nil
	}
	return ctx.GetImports()
}

func GetContext(i types.Instance) types.Context {
	if c, ok := i.GetUserdata().(types.Context); ok {
		return c
	}
	return nil
}

func CopyHostDataToWasm(i types.Instance, data []byte, dataaddrptr, datasizeptr int32) error {
	addr, err := i.Malloc(int32(len(data)))
	if err != nil {
		return err
	}

	err = i.PutMemory(addr, int32(len(data)), data)
	if err != nil {
		return err
	}

	err = i.PutUint32(dataaddrptr, uint32(addr))
	if err != nil {
		return err
	}

	err = i.PutUint32(datasizeptr, uint32(len(data)))
	if err != nil {
		return err
	}

	return nil
}

func ReadStringFromAddr(i types.Instance, addr int32) (string, error) {
	if addr < 4 {
		return "", consts.RESULT__INVALID_MEM_ACCESS
	}
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()

	size, err := i.GetUint32(addr - 4)
	if err != nil {
		return "", err
	}
	mem, err := i.GetMemory(addr, int32(size))
	if err != nil {
		return "", err
	}
	bytes, err := dec.Bytes(mem)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
