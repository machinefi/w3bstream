package vm

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/machinefi/w3bstream/pkg/types"
)

// TODO change this func to NewInstanceByCode
func NewInstance(ctx context.Context, path string, id types.SFID) error {
	return NewInstanceWithState(ctx, path, id, enums.INSTANCE_STATE__CREATED)
}

func NewInstanceWithState(ctx context.Context, path string, id types.SFID, state enums.InstanceState) error {
	l := types.MustLoggerFromContext(ctx)
	fileSystem := types.MustFileSystemOpFromContext(ctx)

	_, l = l.Start(ctx, "NewInstanceWithState")
	defer l.End()

	code, err := fileSystem.Read(path)
	if err != nil {
		l.Error(err)
		return err
	}

	i, err := wasmtime.NewInstanceByCode(ctx, id, code, state)
	if err != nil {
		l.Error(err)
		return err
	}

	AddInstanceByID(ctx, id, i)
	return nil
}
