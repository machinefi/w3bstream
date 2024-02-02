package vm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewInstance(ctx context.Context, code []byte, id types.SFID, state enums.InstanceState) error {
	ctx, l := logr.Start(ctx, "vm.NewInstance")
	defer l.End()

	ins, err := wasmtime.NewInstanceByCode(ctx, id, code, state)
	if err != nil {
		l.Error(err)
		return err
	}
	AddInstanceByID(ctx, id, ins)
	l.WithValues("instance_id", id, "state", state.String()).Info("instance created")
	return nil
}
