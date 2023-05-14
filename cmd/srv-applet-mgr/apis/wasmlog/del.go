package wasmlog

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/wasmlog"
	"github.com/machinefi/w3bstream/pkg/types"
)

type RemoveWasmLogByInstanceID struct {
	httpx.MethodDelete
	InstanceID types.SFID `in:"path" name:"instanceID"`
}

func (r *RemoveWasmLogByInstanceID) Path() string { return "/:instanceID" }

func (r *RemoveWasmLogByInstanceID) Output(ctx context.Context) (interface{}, error) {
	ctx, err := middleware.MustCurrentAccountFromContext(ctx).
		WithInstanceContextBySFID(ctx, r.InstanceID)
	if err != nil {
		return nil, err
	}
	return nil, wasmlog.Remove(ctx, &wasmlog.CondArgs{InstanceID: r.InstanceID})
}
