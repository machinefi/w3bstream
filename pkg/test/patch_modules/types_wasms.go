package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func TypesWasmInitConfiguration(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		wasm.InitConfiguration,
		func(_ context.Context, _ wasm.Configuration) error { return err },
	)
}

func TypesWasmNewConfigurationByType(patch *gomonkey.Patches, c wasm.Configuration, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		wasm.NewConfigurationByType,
		func(_ enums.ConfigType) (wasm.Configuration, error) { return c, err },
	)
}

func TypesWasmUninitConfiguration(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		wasm.UninitConfiguration,
		func(_ context.Context, _ wasm.Configuration) error { return err },
	)
}
