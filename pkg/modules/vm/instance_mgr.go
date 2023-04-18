package vm

import (
	"context"

	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

var instances = mapx.New[types.SFID, wasm.Instance]()

var (
	ErrNotFound = errors.New("instance not found")
)

func AddInstance(ctx context.Context, i wasm.Instance) types.SFID {
	id := confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID()
	AddInstanceByID(ctx, id, i)
	return id
}

func AddInstanceByID(ctx context.Context, id types.SFID, i wasm.Instance) {
	instances.Store(id, i)
}

func DelInstance(ctx context.Context, id types.SFID) error {
	i, _ := instances.LoadAndRemove(id)
	if i == nil {
		return nil
	}
	return i.Stop(ctx)
}

func StartInstance(ctx context.Context, id types.SFID) error {
	i, ok := instances.Load(id)
	if !ok {
		return ErrNotFound
	}

	if i.State() == enums.INSTANCE_STATE__STARTED {
		return nil
	}

	if err := i.Start(ctx); err != nil {
		return err
	}
	return nil
}

func StopInstance(ctx context.Context, id types.SFID) error {
	i, ok := instances.Load(id)
	if !ok {
		return ErrNotFound
	}
	if err := i.Stop(ctx); err != nil {
		return err
	}
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
