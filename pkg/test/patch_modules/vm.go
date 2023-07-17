package patch_modules

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
)

func VmNewInstance(patch *gomonkey.Patches, err error) *gomonkey.Patches {
	return patch.ApplyFunc(
		vm.NewInstance,
		func(_ context.Context, _ []byte, _ types.SFID, _ enums.InstanceState) error { return err },
	)
}
