package vm

import (
	"context"

	confid "github.com/iotexproject/Bumblebee/conf/id"
	"github.com/iotexproject/Bumblebee/x/mapx"
	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/pkg/types"

	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

var instances = mapx.New[types.SFID, wasm.Instance]()

var (
	ErrNotFound = errors.New("instance not found")
)

func AddInstance(ctx context.Context, i wasm.Instance) types.SFID {
	l := types.MustLoggerFromContext(ctx)
	idg := confid.MustSFIDGeneratorFromContext(ctx)

	_, l = l.Start(ctx, "AddInstance")
	defer l.End()

	id := idg.MustGenSFID()
	instances.Store(id, i)

	l.WithValues("instance", id).Info("created")
	return id
}
func AddInstanceByID(ctx context.Context, id types.SFID, i wasm.Instance) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "AddInstanceByID")
	defer l.End()

	instances.Store(id, i)
	l.WithValues("instance", id).Info("created")
}

func DelInstance(ctx context.Context, id types.SFID) error {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "DelInstance")
	defer l.End()

	i, _ := instances.LoadAndRemove(id)
	if i == nil {
		return ErrNotFound
	}
	if i.State() == enums.INSTANCE_STATE__STARTED {
		i.Stop()
	}
	return nil
}

func StartInstance(ctx context.Context, id types.SFID) error {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "StartInstance")
	defer l.End()

	l = l.WithValues("instance", id)

	i, ok := instances.Load(id)
	if !ok {
		l.Error(ErrNotFound)
		return ErrNotFound
	}

	go i.Start(ctx)

	return nil
}

func StopInstance(ctx context.Context, id types.SFID) error {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "StopInstance")
	defer l.End()

	l = l.WithValues("instance", id)

	i, ok := instances.Load(id)
	if !ok {
		l.Warn(ErrNotFound)
		return ErrNotFound
	}
	i.Stop()
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

func GetConsumer(id types.SFID) wasm.EventConsumer {
	i, ok := instances.Load(id)
	if !ok || i == nil {
		return nil
	}
	return i.(wasm.EventConsumer)
}
