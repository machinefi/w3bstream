package vm

import (
	"context"
	"os"

	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewInstance(ctx context.Context, path string, id types.SFID) error {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "NewInstanceWithID")
	defer l.End()

	code, err := os.ReadFile(path)
	if err != nil {
		l.Error(err)
		return err
	}

	i, err := wasmtime.NewInstanceByCode(ctx, id, code)
	if err != nil {
		l.Error(err)
		return err
	}

	AddInstanceByID(ctx, id, i)
	return nil
}
