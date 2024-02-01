package wasmtime

import (
	"context"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	abitypes "github.com/machinefi/w3bstream/pkg/modules/wasm/abi/types"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/host"
	"github.com/machinefi/w3bstream/pkg/modules/wasm/runtime"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

const Engine = "wasmtime"

type Instance struct {
	instance abitypes.Instance
	abictx   abitypes.Context
}

func NewInstanceByCode(ctx context.Context, id types.SFID, code []byte, st enums.InstanceState) (i *Instance, err error) {
	ctx, l := logr.Start(ctx, "vm.NewInstanceByCode")
	defer l.End()

	instance, err := runtime.NewRuntime(ctx, Engine, id.String(), code)
	if err != nil {
		return nil, err
	}
	abictx := host.GetContext(instance)
	if abictx == nil {
		return nil, errors.New("failed to get abi context")
	}

	if st == enums.INSTANCE_STATE__STARTED {
		if err = instance.Start(); err != nil {
			return nil, err
		}
	}

	return &Instance{
		instance: instance,
		abictx:   abictx,
	}, nil
}

func (i *Instance) ID() string {
	return i.instance.ID()
}

func (i *Instance) Start(ctx context.Context) error {
	ctx, l := logr.Start(ctx, "vm.Instance.Start", "instance_id", i.ID())
	defer l.End()

	return i.instance.Start()
}

func (i *Instance) Stop(ctx context.Context) error {
	ctx, l := logr.Start(ctx, "vm.Instance.Stop", "instance_id", i.ID())
	defer l.End()

	i.instance.Stop()
	return nil
}

func (i *Instance) State() wasm.InstanceState {
	if i.instance.Started() {
		return enums.INSTANCE_STATE__STARTED
	}
	return enums.INSTANCE_STATE__STOPPED
}

func (i *Instance) HandleEvent(ctx context.Context, fn, eventType string, data []byte) *wasm.EventHandleResult {
	ctx, l := logr.Start(ctx, "vm.Instance.HandleEvent")
	defer l.End()

	if !i.instance.Acquire() {
		return &wasm.EventHandleResult{
			InstanceID: i.ID(),
			Code:       wasm.ResultStatusCode_Failed,
			ErrMsg:     status.InstanceNotRunning.Key() + "_acquire",
		}
	}
	defer i.instance.Release()

	res, fuel, err := i.abictx.GetExports().OnEventReceived(fn, eventType, data)
	if err != nil {
		l.Error(err)
		return &wasm.EventHandleResult{
			InstanceID: i.ID(),
			Code:       wasm.ResultStatusCode_Failed,
			ErrMsg:     err.Error(),
		}
	}
	code, _ := res.(int32)
	l.WithValues("code", code, "consumed", fuel).Info("")
	return &wasm.EventHandleResult{
		InstanceID: i.ID(),
		Code:       wasm.ResultStatusCode(code),
	}
}
