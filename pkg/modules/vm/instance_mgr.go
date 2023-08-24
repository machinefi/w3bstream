package vm

import (
	"context"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

var instances = mapx.New[types.SFID, wasm.Instance]()

var (
	ErrNotFound = errors.New("instance not found")
)

func AddInstanceByID(ctx context.Context, id types.SFID, i wasm.Instance) {
	ctx, l := logr.Start(ctx, "modules.vm.AddInstanceByID")
	defer l.End()

	instances.Store(id, i)
	l.WithValues("instance", id).Info("created")
}

func DelInstance(ctx context.Context, id types.SFID) error {
	ctx, l := logr.Start(ctx, "modules.vm.DelInstance")
	defer l.End()

	i, _ := instances.LoadAndRemove(id)
	if i == nil {
		return ErrNotFound
	}
	return i.Stop(ctx)
}

func StartInstance(ctx context.Context, id types.SFID) error {
	ctx, l := logr.Start(ctx, "modules.vm.StartInstance")
	defer l.End()

	l = l.WithValues("instance", id)

	i, ok := instances.Load(id)
	if !ok {
		l.Error(ErrNotFound)
		return ErrNotFound
	}

	if i.State() == enums.INSTANCE_STATE__STARTED {
		return nil
	}

	if err := i.Start(ctx); err != nil {
		l.Error(err)
		return err
	}
	l.Info("started")
	return nil
}

func StopInstance(ctx context.Context, id types.SFID) error {
	ctx, l := logr.Start(ctx, "modules.vm.StopInstance")
	defer l.End()

	l = l.WithValues("instance", id)

	i, ok := instances.Load(id)
	if !ok {
		l.Warn(ErrNotFound)
		return ErrNotFound
	}
	if err := i.Stop(ctx); err != nil {
		l.Error(err)
		return err
	}
	l.Info("stopped")
	return nil
}

func GetInstanceState(id types.SFID) (enums.InstanceState, bool) {
	i, ok := instances.Load(id)
	if !ok {
		return enums.INSTANCE_STATE_UNKNOWN, false
	}
	return i.State(), true
}

func GetConsumer(id types.SFID) wasm.Instance {
	i, ok := instances.Load(id)
	if !ok || i == nil {
		return nil
	}
	return i
}
