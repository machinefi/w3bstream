package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func TypesWasmInitConfiguration(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		wasm.InitConfiguration,
		func(_ context.Context, _ wasm.Configuration) error { return err },
	)
}
