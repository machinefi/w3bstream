package vm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewInstance(ctx context.Context, code []byte, id types.SFID, state enums.InstanceState) error {
	ctx, l := logr.Start(ctx, "modules.vm.NewInstance")
	defer l.End()

	ins, err := wasmtime.NewInstanceByCode(ctx, id, code, state)
	if err != nil {
		return err
	}
	AddInstanceByID(ctx, id, ins)
	return nil
}
