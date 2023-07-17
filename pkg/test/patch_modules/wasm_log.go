package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/modules/wasmlog"
)

func WasmLogRemove(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		wasmlog.Remove,
		func(_ context.Context, _ *wasmlog.CondArgs) error { return err },
	)
}
