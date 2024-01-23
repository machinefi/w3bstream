package internal

import (
	"golang.org/x/text/encoding/unicode"

	"github.com/machinefi/w3bstream/pkg/modules/wasm"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/consts"
)

func GetImportsHandler(i Instance) ImportsHandler {
	ctx := GetContext(i)
	if ctx == nil {
		return nil
	}
	return ctx.GetImports()
}

func GetContext(i Instance) wasm.Context {
	if c, ok := i.GetUserdata().(wasm.Context); ok {
		return c
	}
	return nil
}

func CopyDataToInstance(i Instance, data []byte, dataaddrptr, datasizeptr int32) error {
	addr, err := i.Malloc(int32(len(data)))
	if err != nil {
		return err
	}

	err = i.PutMemory(addr, uint64(len(data)), data)
	if err != nil {
		return err
	}

	err = i.PutUint32(uint64(dataaddrptr), uint32(addr))
	if err != nil {
		return err
	}

	err = i.PutUint32(uint64(datasizeptr), uint32(len(data)))
	if err != nil {
		return err
	}

	return nil
}

func ReadStringFromAddr(i Instance, addr int32) (string, error) {
	if addr < 4 {
		return "", consts.RESULT__INVALID_MEM_ACCESS
	}
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()

	size, err := i.GetUint32(uint64(addr - 4))
	if err != nil {
		return "", err
	}
	mem, err := i.GetMemory(uint64(addr), uint64(size))
	if err != nil {
		return "", err
	}
	bytes, err := dec.Bytes(mem)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
